package plugin

import (
	"net/http"
	"regexp"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

// Router is a task which maps paths to routes
type Router interface {
	Task
	http.Handler

	// Add a prefix/path mapping to a handler for one or more HTTP methods
	AddHandler(Gateway, *regexp.Regexp, http.HandlerFunc, ...string) error

	// Add middleware handler to the router given unique name
	AddMiddleware(string, func(http.HandlerFunc) http.HandlerFunc) error
}
