package main

import (
	// Modules
	httpserver "github.com/mutablelogic/terraform-provider-nginx/pkg/httpserver"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

func Config() TaskPlugin {
	return httpserver.Config{}
}
