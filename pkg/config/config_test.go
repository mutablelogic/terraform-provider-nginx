package config_test

import (
	"os"
	"testing"

	// Module import
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/config"
)

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	TEST_DIR = "../../etc/test"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Config_001(t *testing.T) {
	// Create an "enabled" folder
	enabled_path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(enabled_path)

	// Create a configuration
	cfg, err := Config{
		AvailablePath: TEST_DIR,
		EnabledPath:   enabled_path,
	}.New()

	// Check for errors
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cfg)
	}
}

func Test_Config_002(t *testing.T) {
	// Create an "enabled" folder
	enabled_path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(enabled_path)

	// Create a configuration
	cfg, err := Config{
		AvailablePath: TEST_DIR,
		EnabledPath:   enabled_path,
	}.New()

	// Check for errors
	if err != nil {
		t.Error(err)
	}

	// Enumerate available files
	files, err := cfg.EnumerateAvailable()
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		t.Log(file)
	}
}

func Test_Config_003(t *testing.T) {
	// Create an "enabled" folder
	enabled_path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(enabled_path)

	// Create a configuration
	cfg, err := Config{
		AvailablePath: TEST_DIR,
		EnabledPath:   enabled_path,
	}.New()

	// Check for errors
	if err != nil {
		t.Error(err)
	}

	// Create new file
	file, err := cfg.Create("test.conf", []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	// Link file to enabled folder
	if link, err := cfg.Link(file); err != nil {
		t.Fatal(err)
	} else {
		t.Log("Linked: ", file, "=>", link)
	}

	// Check enabled
	if enabled, err := cfg.Enabled(file); err != nil {
		t.Fatal(err)
	} else if enabled == false {
		t.Error("File is not enabled", file)
	}

	// Remove file
	if err := cfg.Remove(file); err != nil {
		t.Fatal(err)
	}
}

func Test_Config_004(t *testing.T) {
	// Create an "enabled" folder
	enabled_path, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(enabled_path)

	// Create a configuration
	cfg, err := Config{
		AvailablePath: TEST_DIR,
		EnabledPath:   enabled_path,
	}.New()

	// Check for errors
	if err != nil {
		t.Error(err)
	}

	// Enumerate available files
	files, err := cfg.EnumerateAvailable()
	if err != nil {
		t.Fatal(err)
	}

	// Check enabled
	for _, file := range files {
		if enabled, err := cfg.Enabled(file); err != nil {
			t.Fatal(err)
		} else if enabled == true {
			t.Error("Expected unenabled file", file)
		}
	}

	// Link
	for _, file := range files {
		if link, err := cfg.Link(file); err != nil {
			t.Fatal(err)
		} else {
			t.Log("Linked: ", file, "=>", link)
		}
	}

	// Check enabled
	for _, file := range files {
		if enabled, err := cfg.Enabled(file); err != nil {
			t.Fatal(err)
		} else if enabled == false {
			t.Error("Expected enabled file", file)
		}
	}

	// Unlink
	for _, file := range files {
		if err := cfg.Unlink(file); err != nil {
			t.Fatal(err)
		} else {
			t.Log("Unlinked: ", file)
		}
	}

	// Enumerable enabled files
	files2, err := cfg.EnumerateEnabled()
	if err != nil {
		t.Fatal(err)
	}
	if len(files2) > 0 {
		t.Error("Expected no enabled files")
	}
}
