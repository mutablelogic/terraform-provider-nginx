package util_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	// Module import
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Writable_001(t *testing.T) {
	// Create a temprary directory
	dir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create a folder which is writable
	w := filepath.Join(dir, "writable")
	if err := os.Mkdir(w, 0755); err != nil {
		t.Fatal(err)
	} else if writable, err := IsWritableDir(w); err != nil {
		t.Fatal(err)
	} else if !writable {
		t.Error("Folder is not writable", w)
	}

	// Create a folder which is readonly
	ro := filepath.Join(dir, "readonly")
	if err := os.Mkdir(ro, 0555); err != nil {
		t.Fatal(err)
	} else if writable, err := IsWritableDir(ro); err != nil {
		t.Fatal(err)
	} else if writable {
		t.Error("Folder is writable, but should be readonly", ro)
	}
}
