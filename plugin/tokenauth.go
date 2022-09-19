package plugin

import (
	"time"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

// TokenAuth stores tokens for authentication
type TokenAuth interface {
	Task

	// Return true if a token associated with the name already exists
	Exists(string) bool

	// Create a new token associated with a name and return it.
	Create(string) (string, error)

	// Revoke a token associated with a name. For the admin token, it is
	// rotated rather than revoked.
	Revoke(string) error

	// Return all token names and their last access times
	Enumerate() map[string]time.Time

	// Returns the name of the token if a value matches. Updates
	// the access time for the token. If token with value not
	// found, then return empty string
	Matches(string) string
}
