package main

import (
	// Modules
	gateway "github.com/mutablelogic/terraform-provider-nginx/pkg/mdns-gateway"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

func Config() TaskPlugin {
	return gateway.Config{}
}
