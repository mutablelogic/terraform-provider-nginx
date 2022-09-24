package nginx_test

import (
	"context"
	"testing"

	// Module import
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"

	//plugin "github.com/mutablelogic/terraform-provider-nginx/plugin"

	// Namespace imports
	//. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/nginx"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Router_001(t *testing.T) {
	// Create a provider, register http server and router
	p := provider.New()
	nginx, err := p.New(context.Background(), Config{
		Path: "/etc/nginx",
	})
	if err != nil {
		t.Fatal(err)
	}
	if nginx == nil {
		t.Fatal("Unexpected nil returned from NewRouter")
	} else {
		t.Log(nginx)
	}
}
