package nginx_test

import (
	"context"
	"os"
	"testing"

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
	if configs, err := nginx.(plugin.Nginx).Enumerate(); err != nil {
		t.Error(err)
	} else {
		t.Log(configs)
	}

}
