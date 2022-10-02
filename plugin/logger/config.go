package main

import (
	// Modules
	logger "github.com/mutablelogic/terraform-provider-nginx/pkg/logger"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

func Config() TaskPlugin {
	return logger.Config{}
}
