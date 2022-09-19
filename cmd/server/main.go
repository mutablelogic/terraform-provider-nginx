package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// Moudule imports
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"
)

var (
	flagAddr = flag.String("addr", ":8080", "Address to listen on")
)

func main() {
	flag.Parse()

	provider := provider.New()
	ctx := HandleSignal()

	// Built-in tokenauth
	if _, err := provider.New(ctx, tokenauth.Config{Delta: time.Second}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Run the server and plugins
	fmt.Println("Running server, press CTRL-C to exit")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := provider.Run(HandleSignal()); err != nil {
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
