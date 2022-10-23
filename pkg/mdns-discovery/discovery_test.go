package mdns_discovery_test

import (
	"context"
	"sync"
	"testing"
	"time"

	// Modules

	discovery "github.com/mutablelogic/terraform-provider-nginx/pkg/mdns-discovery"
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
)

func Test_Discovery_001(t *testing.T) {
	provider := provider.New()
	_, err := provider.New(context.Background(), discovery.Config{})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(provider)
	}
}

func Test_Discovery_002(t *testing.T) {
	provider := provider.New()
	mdns, err := provider.New(context.Background(), discovery.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Discover for 2 seconds in the background
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Discover services
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if services, err := mdns.(discovery.MDNSTask).Discover(ctx); err != nil {
			t.Error(err)
		} else {
			t.Log("Services=", services)
		}
	}()

	// Run the tasks
	if err := provider.Run(ctx); err != nil {
		t.Fatal(err)
	}

	// Wait until all tasks have finished
	wg.Wait()
}
