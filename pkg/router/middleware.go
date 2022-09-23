package router

import (
	"net/http"
	"sort"

	// Namespace imports
	. "github.com/djthorpe/go-errors"

	// Module imports
	"github.com/mutablelogic/terraform-provider-nginx/pkg/util"
	//. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type middleware struct {
	handlers map[string]func(http.HandlerFunc) http.HandlerFunc
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// AddHandler adds a handler to the router, for a specific prefix and http methods supported.
// If the path argument is nil, then any path under the prefix will match. If the path contains
// a regular expression, then a match is made and any matched parameters of the regular
// expression can be retrieved from the request context.
func (m *middleware) AddMiddleware(name string, fn func(http.HandlerFunc) http.HandlerFunc) error {
	if !util.IsIdentifier(name) || fn == nil {
		return ErrBadParameter.With(name)
	}
	if m.handlers == nil {
		m.handlers = make(map[string]func(http.HandlerFunc) http.HandlerFunc)
	}

	return ErrNotImplemented.With("middleware: ", name)
}

// Wrap a handler with middleware, called from right to left
func (m *middleware) Wrap(fn http.HandlerFunc, middleware ...string) (http.HandlerFunc, error) {
	sort.Reverse(sort.StringSlice(middleware))
	for _, name := range middleware {
		if wrapped, ok := m.handlers[name]; ok {
			fn = wrapped(fn)
		} else {
			return nil, ErrNotFound.With("middleware: ", name)
		}
	}
	return fn, nil
}

/*

// AddMiddleware adds a middleware handler with a unique key.
func (r *router) AddMiddleware(key string, fn func(http.HandlerFunc) http.HandlerFunc) error {
	// Preconditions
	if !reValidName.MatchString(key) {
		return ErrBadParameter.Withf("AddMiddleWare: %q", key)
	}
	if fn == nil {
		return ErrBadParameter.Withf("AddMiddleWare: %q", key)
	}

	// Check for duplicate entry
	r.RLock()
	_, exists := r.middleware[key]
	r.RUnlock()
	if exists {
		return ErrDuplicateEntry.Withf("AddMiddleWare: %q", key)
	}

	// Set middleware mapping
	r.Lock()
	r.middleware[key] = fn
	r.Unlock()

	// Return success
	return nil
}

// SetMiddleware binds an array of middleware functions to a prefix. The prefix should
// already exist in the router.
func (r *router) SetMiddleware(prefix string, chain ...string) error {
	prefix = normalizePath(prefix, true)
	fmt.Println("SetMiddleware", prefix, chain)

	return nil
}
*/
