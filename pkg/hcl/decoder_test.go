package hcl_test

import (
	"os"
	"path/filepath"
	"testing"

	// Module import
	httpserver "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/hcl"
)

func Test_Decoder_001(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a HCL decoder and register httpserver and router
	decoder := NewDecoder()
	decoder.Register(httpserver.ServerConfig{})
	decoder.Register(httpserver.RouterConfig{})

	// Parse HCL to get the plugins
	plugins, err := decoder.Parse(os.DirFS("/"), filepath.Join(wd, HCL_TESTS))
	if err != nil {
		t.Fatal(err)
	}

	// Describe the plugins. References to variables and tasks have not yet been resolved
	for _, plugin := range plugins {
		t.Log(plugin.Name(), "=>", plugin)
	}
}
