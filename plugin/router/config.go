package main

import (
	// Modules
	router "github.com/mutablelogic/terraform-provider-nginx/pkg/router"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

func Config() TaskPlugin {
	return router.Config{}
}
