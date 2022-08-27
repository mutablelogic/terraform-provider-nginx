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
	auth "github.com/mutablelogic/terraform-provider-nginx/pkg/auth"
	authgw "github.com/mutablelogic/terraform-provider-nginx/pkg/authgw"
	httpserver "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"
	kernel "github.com/mutablelogic/terraform-provider-nginx/pkg/kernel"
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
	auth, err := auth.Config{Path: *flagPath}.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else if err := kernel.Add("auth", auth); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else if err := kernel.Add("authgw", authgw.New(auth, "/v1/auth")); err != nil {
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

	// AuthAdmin middleware makes sure only admin authenticated users can access the authgw endpoint
	//kernel.AddMiddleware("authgw", authgw.AuthAdmin)

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
