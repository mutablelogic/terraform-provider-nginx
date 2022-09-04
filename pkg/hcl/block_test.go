package hcl_test

import (
	"reflect"
	"testing"

	// Module import
	httpserver "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/hcl"
)

func Test_Block_001(t *testing.T) {
	if spec, err := SpecForType(reflect.TypeOf(httpserver.ServerConfig{}), TagName, "server"); err != nil {
		t.Fatal(err)
	} else {
		t.Log(spec)
	}
}
