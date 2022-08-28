package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// Moudule imports
	authgw "github.com/mutablelogic/terraform-provider-nginx/pkg/authgw"
	httpserver "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"
	kernel "github.com/mutablelogic/terraform-provider-nginx/pkg/kernel"
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"
)

var (
	flagPath = flag.String("path", "/var/lib/nginx-gw", "Path to application storage")
	flagAddr = flag.String("addr", "", "Server listening address or socket path")
)

func main() {
	flag.Parse()

	// Create a new kernel
	kernel := kernel.New()

	// Auth task
	auth, err := tokenauth.Config{Path: *flagPath}.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else if err := kernel.Add("tokenauth", auth); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else if err := kernel.Add("authgw", authgw.New(auth, "/auth/v1")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// Router task
	if err := kernel.Add("router", httpserver.NewRouter()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// Server task
	server, err := httpserver.Config{Addr: *flagAddr}.New(kernel.Get("router").(http.Handler))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else if err := kernel.Add("server", server); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// Only admin authenticated users should be able to access pages under the /auth/v1 endpoint
	if err := kernel.SetMiddleware("/auth/v1", "token-auth", "token-admin-auth"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// Go routine for receiving events
	go func() {
		for evt := range kernel.C() {
			fmt.Println(evt)
		}
	}()

	// Run the kernel
	fmt.Fprintln(os.Stderr, "Press CTRL-C to exit")
	if err := kernel.Run(HandleSignal()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	fmt.Fprintln(os.Stderr, "\nCompleted successfully")
}

func HandleSignal() context.Context {
	// Handle signals - call cancel when interrupt received
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
	}()
	return ctx
}

/*
task "auth-gw" {
	path = "/var/lib/nginx-gw"
    prefix = "/auth/v1"
    middleware = [ "token-auth", "token-admin-auth" ]
}
task "nginx-gw" {
    path = "/var/lib/nginx-gw"
}
*/
