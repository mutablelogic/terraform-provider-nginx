package plugin

import (
	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

// Nginx provides management of configurations
type Nginx interface {
	Task

	// Return all configurations
	Enumerate() ([]NginxConfig, error)
}

type NginxConfig interface {
	// Return the name of the configuration
	Name() string

	// Return the state of the configuration
	Enabled() bool
}
