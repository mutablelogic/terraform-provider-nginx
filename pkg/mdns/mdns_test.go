package mdns_test

import (
	"context"
	"sync"
	"testing"
	"time"

	// Modules
	mdns "github.com/mutablelogic/terraform-provider-nginx/pkg/mdns"
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
)

func Test_MDNS_001(t *testing.T) {
	provider := provider.New()
	task, err := provider.New(context.Background(), mdns.Config{})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(task)
	}
}

func Test_MDNS_002(t *testing.T) {
	provider := provider.New()
	_, err := provider.New(context.Background(), mdns.Config{})
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for evt := range provider.Sub() {
			t.Log(evt)
		}
	}()

	if err := provider.Run(ctx); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}
