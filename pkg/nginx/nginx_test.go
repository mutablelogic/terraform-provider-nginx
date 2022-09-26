package nginx_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	// Module import
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
	plugin "github.com/mutablelogic/terraform-provider-nginx/plugin"

	// Namespace imports
	//. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/nginx"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Nginx_001(t *testing.T) {
	// Set up temporary folder for enabled
	tmpdir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create a provider, register http server and router
	p := provider.New()
	nginx, err := p.New(context.Background(), Config{
		Available: "../../etc/test/nginx",
		Enabled:   tmpdir,
	})
	if err != nil {
		t.Fatal(err)
	} else if nginx == nil {
		t.Fatal("Unexpected nil returned from New")
	} else {
		t.Log(nginx)
	}

	// Enumerate files
	configs, err := nginx.(plugin.Nginx).Enumerate()
	if err != nil {
		t.Error(err)
	}

	// Enable all configs
	for _, config := range configs {
		if !config.Enabled() {
			if err := nginx.(plugin.Nginx).Enable(config); err != nil {
				t.Error(err)
			}
		} else {
			t.Error("Unexpected enabled config", config.Name())
		}
	}

	// Disable all configs
	for _, config := range configs {
		if config.Enabled() {
			if err := nginx.(plugin.Nginx).Disable(config); err != nil {
				t.Error(err)
			}
		} else {
			t.Error("Unexpected disabled config", config.Name())
		}
	}
}

func Test_Nginx_002(t *testing.T) {
	// Set up temporary folder for enabled
	tmpdir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create a provider, register http server and router
	p := provider.New()
	nginx, err := p.New(context.Background(), Config{
		Available: "../../etc/test/nginx",
		Enabled:   tmpdir,
	})
	if err != nil {
		t.Fatal(err)
	} else if nginx == nil {
		t.Fatal("Unexpected nil returned from New")
	} else {
		t.Log(nginx)
	}

	// Run for 5 seconds
	ctx, cancel := context.WithCancel(context.Background())

	// Run
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		p.Run(ctx)
	}()

	// Create a new config
	cfg, err := nginx.(plugin.Nginx).Create("test", []byte("new config"))
	if err != nil {
		t.Error(err)
	}

	// Enable config
	if err := nginx.(plugin.Nginx).Enable(cfg); err != nil {
		t.Error(err)
	}

	// Revoke config
	if err := nginx.(plugin.Nginx).Revoke(cfg); err != nil {
		t.Error(err)
	}

	// Wait for five seconds
	time.Sleep(5 * time.Second)

	// Cancel and wait
	cancel()
	wg.Wait()
}
