package router_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	// Module import
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/router"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Router_001(t *testing.T) {
	// Create a provider, register http server and router
	p := provider.New()
	router, err := p.New(context.Background(), Config{})
	if err != nil {
		t.Fatal(err)
	}
	if router == nil {
		t.Fatal("Unexpected nil returned from NewRouter")
	} else {
		t.Log(router)
	}
}

func Test_Router_002(t *testing.T) {
	// Create a provider, register http server and router
	p := provider.New()
	router, err := p.New(context.Background(), Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Add a route for '/'
	if err := router.(Router).AddHandler("/", nil, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/"))
	}); err != nil {
		t.Error(err)
	}
	// Add a route for '/A'
	if err := router.(Router).AddHandler("/A", nil, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/A"))
	}); err != nil {
		t.Error(err)
	}
	// Add a route for '/AA'
	if err := router.(Router).AddHandler("/AA", nil, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/AA"))
	}); err != nil {
		t.Error(err)
	}
	// Add a route for '/' with regexp
	if err := router.(Router).AddHandler("/", regexp.MustCompile("^/(test1)"), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/test1"))
	}); err != nil {
		t.Error(err)
	}
	// Add a route for '/AA' with regexp
	if err := router.(Router).AddHandler("/AA", regexp.MustCompile("^/(test2)"), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/AA/test2"))
	}); err != nil {
		t.Error(err)
	}
	// Add a route for '/AA' with regexp
	if err := router.(Router).AddHandler("/AA", regexp.MustCompile("^/(test3)"), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/AA/test3"))
	}); err != nil {
		t.Error(err)
	}

	tests := []struct {
		Method, Path string
		Code         int
		Expected     string
	}{
		{http.MethodGet, "/", http.StatusOK, "/"},
		{http.MethodGet, "/test", http.StatusOK, "/"},
		{http.MethodPost, "/test", http.StatusMethodNotAllowed, ""},
		{http.MethodGet, "/A", http.StatusOK, "/"},
		{http.MethodGet, "/A/test1", http.StatusOK, "/A"},
		{http.MethodGet, "/AB", http.StatusOK, "/"},
		{http.MethodGet, "/AA/test2", http.StatusOK, "/AA/test2"},
		{http.MethodGet, "/AA/test3", http.StatusOK, "/AA/test3"},
		{http.MethodGet, "/AAA/test3", http.StatusOK, "/"},
		{http.MethodGet, "/AAtest3", http.StatusOK, "/"},
	}

	for i, test := range tests {
		w := httptest.NewRecorder()
		router.(http.Handler).ServeHTTP(w, httptest.NewRequest(test.Method, test.Path, nil))
		if status := w.Result().StatusCode; status != test.Code {
			t.Error("Test", i, ": unexpected status code: ", status)
		} else if body, _ := io.ReadAll(w.Result().Body); test.Expected != "" && string(body) != test.Expected {
			t.Errorf("Test %d: unexpected body: %q", i, body)
		}
	}
}
