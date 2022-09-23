package tokenauth_gateway_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// Module imports
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
	router "github.com/mutablelogic/terraform-provider-nginx/pkg/router"
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"
	gateway "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth-gateway"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

func Test_TokenAuthGateway_001(t *testing.T) {
	provider := provider.New()
	ctx := context.Background()

	// Create an "tokens" folder
	path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	tokenauth, err := provider.New(ctx, tokenauth.Config{Path: path})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(tokenauth)
	}
	router, err := provider.New(ctx, router.Config{})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(router)
	}
	gateway, err := provider.New(ctx, gateway.Config{Auth: tokenauth, Router: router})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(gateway)
	}
}

func Test_TokenAuthGateway_002(t *testing.T) {
	provider := provider.New()
	ctx := context.Background()

	// Create an "tokens" folder
	path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	tokenauth, err := provider.New(ctx, tokenauth.Config{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	router, err := provider.New(ctx, router.Config{})
	if err != nil {
		t.Fatal(err)
	}
	gateway, err := provider.New(ctx, gateway.Config{Auth: tokenauth, Router: router})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(gateway)
	}

	// Check /list method
	w := httptest.NewRecorder()
	router.(http.Handler).ServeHTTP(w, httptest.NewRequest(http.MethodGet, gateway.(Gateway).Prefix()+"/", nil))
	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Error("unexpected status code: ", status)
	} else {
		body, _ := io.ReadAll(w.Result().Body)
		t.Log(string(body))
	}
}
