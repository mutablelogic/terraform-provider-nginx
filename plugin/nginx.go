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

	// Enable a configuration
	Enable(NginxConfig) error

	// Revoke a configuration
	Revoke(NginxConfig) error
}

// NginxConfig provides a configuration that can be enabled or revoked
type NginxConfig interface {
	// Return the name of the configuration
	Name() string

	// Return the state of the configuration
	Enabled() bool
}
