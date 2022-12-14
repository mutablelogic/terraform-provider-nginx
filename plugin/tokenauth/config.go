package main

import (
	// Modules
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

func Config() TaskPlugin {
	return tokenauth.Config{}
}
