package httpserver_test

import (
	"context"
	"testing"

	// Module import
	"github.com/mutablelogic/terraform-provider-nginx/pkg/provider"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Server_001(t *testing.T) {
	// Create a provider, register http server and router
	p := provider.New()
	if err := p.Register(ServerConfig{}); err != nil {
		t.Fatal(err)
	}
	if err := p.Register(RouterConfig{}); err != nil {
		t.Fatal(err)
	}

	// Create a router
	router, err := p.New(context.Background(), RouterConfig{})
	if err != nil {
		t.Fatal(err)
	}

	// Create a server
	if server, err := p.New(context.Background(), ServerConfig{
		Label:  "label",
		Router: router,
	}); err != nil {
		t.Fatal(err)
	} else {
		t.Log(server)
	}
}
