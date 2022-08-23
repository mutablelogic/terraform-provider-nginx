package auth_test

import (
	"os"
	"testing"

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

	// Check for errors
	if err != nil {
		t.Error(err)
	} else {
		t.Log(auth)
	}
}
