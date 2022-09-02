package hcl_test

import (
	"os"
	"testing"

	// Module import
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/hcl"
)

const (
	HCL_TESTS = "../../etc/test/hcl"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_HCL_001(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a hcl parser
	hcl, err := New(wd, HCL_TESTS, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if hcl == nil {
		t.Fatal("Unexpected nil returned from New")
	}
	t.Log(hcl)
}

func Test_HCL_002(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Create a hcl parser
	hcl, err := New(wd, HCL_TESTS, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	if err := hcl.Parse(); err != nil {
		t.Fatal(err)
	}
}
