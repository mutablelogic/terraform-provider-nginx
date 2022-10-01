package provider_test

import (
	"context"
	"sync"
	"testing"
	"time"

	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
)

func Test_Task_001(t *testing.T) {
	// create provider
	cfg := provider.Config{Label: t.Name()}
	provider := provider.New()
	if provider == nil {
		t.Fatal("Expected provider")
	}

	// Receive events from provider
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for evt := range provider.C() {
			t.Log("Received", evt)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	task, err := provider.New(ctx, cfg)
	if err != nil {
		t.Fatal("Expected task, got", err)
	} else {
		t.Log("Created", task)
	}

	go func() {
		// Wait for a second
		time.Sleep(time.Second)
		// Cancel provider
		t.Log("Cancel")
		cancel()
	}()

	// Run provider
	t.Log("Running")
	if err := provider.Run(ctx); err != nil {
		t.Fatal("Expected provider to run, got", err)
	}
	t.Log("End run")

	// Wait for events to be closed
	wg.Wait()

}
