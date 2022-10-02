package config_test

import (
	"os"
	"path/filepath"
	"testing"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/config"
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
	if resources, err := ParseJSON(fs, "json/*.json"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(resources)
	}
}
