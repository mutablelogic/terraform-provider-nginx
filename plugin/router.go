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
	AddHandler(prefix string, path *regexp.Regexp, fn http.HandlerFunc, methods ...string) error

	// Add middleware for a unique name
	//AddMiddleware(name string, fn func(http.HandlerFunc) http.HandlerFunc) error

	// Set middleware for a prefix. Called from left to right.
	//SetMiddleware(prefix string, chain ...string) error
}
