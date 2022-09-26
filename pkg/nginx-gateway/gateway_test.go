package nginx_gateway_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	// Module imports
	nginx "github.com/mutablelogic/terraform-provider-nginx/pkg/nginx"
	gateway "github.com/mutablelogic/terraform-provider-nginx/pkg/nginx-gateway"
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
	router "github.com/mutablelogic/terraform-provider-nginx/pkg/router"
	// Namespace imports
	//. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

func Test_NginxGateway_001(t *testing.T) {
	provider := provider.New()
	ctx := context.Background()

	// Create an folder for enabled configurations
	tmpdir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create tasks and add them to the provider
	nginx, err := provider.New(ctx, nginx.Config{
		Available: "../../etc/test/nginx",
		Enabled:   tmpdir,
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(nginx)
	}
	router, err := provider.New(ctx, router.Config{})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(router)
	}
	gateway, err := provider.New(ctx, gateway.Config{Nginx: nginx, Router: router})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(gateway)
	}
}

func Test_NginxGateway_002(t *testing.T) {
	provider := provider.New()
	ctx := context.Background()

	// Create an folder for enabled configurations
	tmpdir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create tasks and add them to the provider
	nginx, err := provider.New(ctx, nginx.Config{
		Available: "../../etc/test/nginx",
		Enabled:   tmpdir,
	})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(nginx)
	}
	router, err := provider.New(ctx, router.Config{})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(router)
	}
	gateway, err := provider.New(ctx, gateway.Config{Nginx: nginx, Router: router})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(gateway)
	}

	// Run tasks until cancel
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	wg.Add(2)
	go func() {
		defer wg.Done()
		t.Log("Running provider")
		if err := provider.Run(ctx); err != nil {
			t.Error(err)
		}
		t.Log("Finished running provider")
	}()
	go func() {
		defer wg.Done()
		t.Log("Running event handler")
		for event := range provider.C() {
			t.Log(event)
		}
		t.Log("Finished event handler")
	}()

	time.Sleep(10 * time.Second)
	cancel()
	wg.Wait()
}
