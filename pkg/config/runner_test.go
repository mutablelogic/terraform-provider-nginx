package config_test

import (
	"context"
	"os"
	"testing"
	"time"

	// Module import
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/config"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Runner_001(t *testing.T) {
	// Create an "enabled" folder
	enabled_path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(enabled_path)

	// Create a configuration
	runner, err := Config{
		AvailablePath: TEST_DIR,
		EnabledPath:   enabled_path,
	}.NewRunner()

	// Check for errors
	if err != nil {
		t.Error(err)
	}

	// Run for 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := runner.Run(ctx); err != nil {
		t.Error(err)
	}
}
