package hcl_test

import (
	"os"
	"path/filepath"
	"testing"

	// Module import

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/hcl"
)

const (
	HCL_TESTS = "../../etc/test/hcl"
)

func Test_Parse_001(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if body, err := Parse(os.DirFS("/"), filepath.Join(wd, HCL_TESTS)); err != nil {
		t.Fatal(err)
	} else {
		t.Log("parsed=", body)
	}
}
