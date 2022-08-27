package auth_test

import (
	"context"
	"os"
	"testing"
	"time"

	// Module import
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/auth"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Auth_001(t *testing.T) {
	// Create an "tokens" folder
	path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	// Create a configuration
	auth, err := Config{Path: path}.New()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(auth)
}

func Test_Auth_002(t *testing.T) {
	// Create an "tokens" folder
	path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	// Create a configuration
	auth, err := Config{Path: path, Delta: time.Second}.New()
	if err != nil {
		t.Fatal(err)
	}

	// Write out events
	go func() {
		for evt := range auth.C() {
			t.Log(evt)
		}
	}()

	// Revoke the admin token once to rotate it
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			t.Log("Revoke admin token")
			if err := auth.Revoke(AdminToken); err != nil {
				t.Error(err)
			}
		}
	}()

	// Run in foreground for a few seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := auth.Run(ctx); err != nil {
		t.Error(err)
	}
}

func Test_Auth_003(t *testing.T) {
	// Create an "tokens" folder
	path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	// Create a configuration
	auth, err := Config{Path: path, Delta: time.Second}.New()
	if err != nil {
		t.Fatal(err)
	}

	// Write out events
	go func() {
		for evt := range auth.C() {
			t.Log(evt)
		}
	}()

	// Create and then revoke a token
	go func() {
		value, err := auth.Create("test")
		if err != nil {
			t.Error(err)
		} else {
			t.Log("token=", value)
		}
		time.Sleep(2 * time.Second)
		t.Log("Testing test token value for match")
		if matches := auth.Matches(value); matches != "test" {
			t.Error("Expected token to be 'test'")
		}
		time.Sleep(2 * time.Second)
		t.Log("Revoking test token")
		if err := auth.Revoke("test"); err != nil {
			t.Error(err)
		} else {
			t.Log("Revoked token=", value)
		}
		time.Sleep(2 * time.Second)
		if matches := auth.Matches(value); matches != "" {
			t.Error("Expected token to be empty")
		}
	}()

	// Run in foreground for a few seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := auth.Run(ctx); err != nil {
		t.Error(err)
	}
}
