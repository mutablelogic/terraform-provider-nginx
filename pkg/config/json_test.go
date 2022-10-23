package config_test

import (
	"os"
	"path/filepath"
	"testing"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/config"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"
)

const (
	baseTestConfigPath = "../../etc/test"
)

func Test_JSON_001(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fs := os.DirFS(filepath.Join(cwd, baseTestConfigPath))
	if resources, err := LoadJSONForPattern(fs, "json/*.json"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(resources)
	}
}

func Test_JSON_002(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fs := os.DirFS(filepath.Join(cwd, baseTestConfigPath))
	if resources, err := LoadJSONForPattern(fs, "json/httpserver.json"); err != nil {
		t.Fatal(err)
	} else if len(resources) != 1 {
		t.Error("Expected one resource")
	} else if plugin, err := ParseJSONResource(fs, resources[0], httpserver.Config{}); err != nil {
		t.Error(err)
	} else {
		t.Log(plugin)
	}
}
