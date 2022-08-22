package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// Moudule imports

	server "github.com/mutablelogic/terraform-provider-nginx/pkg/server"
)

var (
	flagAvailable = flag.String("available", "/etc/nginx/sites-available", "Path to the available sites")
	flagEnabled   = flag.String("enabled", "/etc/nginx/sites-enabled", "Path to the enabled sites")
	flagAddr      = flag.String("addr", ":8080", "Address to listen on")
)

func main() {
	flag.Parse()

	// Create a new gateway
	gateway, err := NewGateway(*flagAvailable, *flagEnabled)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create a server and a router
	server, err := server.Config{Addr: *flagAddr, Router: gateway}.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Run the server and plugins
	fmt.Println("Running server, press CTRL-C to exit")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := gateway.Run(HandleSignal()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Run(HandleSignal()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}()

	// Wait for end
	wg.Wait()
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
