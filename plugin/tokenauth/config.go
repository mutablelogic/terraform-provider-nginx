package main

import (
	// Module imports
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

func Config() TaskPlugin {
	return tokenauth.Config{}
}
