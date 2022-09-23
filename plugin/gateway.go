package plugin

import (
	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

// Gateway provides handlers for a gateway
type Gateway interface {
	Task

	// Return the prefix for this gateway
	Prefix() string
}
