package router

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	// Module imports
	context "github.com/mutablelogic/terraform-provider-nginx/pkg/context"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type middlewarefn func(http.HandlerFunc) http.HandlerFunc

type router struct {
	sync.RWMutex
	label      string
	routes     []route
	cache      map[string]*cached
	middleware map[string]middlewarefn
}

type cached struct {
	index   int
	matched []string
}

type route struct {
	prefix  string
	path    *regexp.Regexp
	fn      http.HandlerFunc
	methods []string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (Router, error) {
	r := new(router)
	r.cache = make(map[string]*cached)
	r.middleware = make(map[string]middlewarefn)
	r.label = c.Label

	// Return success
	return r, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r *router) String() string {
	str := "<router"
	str += fmt.Sprintf(" label=%q", r.label)
	for _, route := range r.routes {
		str += fmt.Sprintf(" %q %q => %q", route.prefix, route.path, route.methods)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// AddHandler adds a handler to the router, for a specific prefix and http methods supported.
// If the path argument is nil, then any path under the prefix will match. If the path contains
// a regular expression, then a match is made and any matched parameters of the regular
// expression can be retrieved from the request context.
func (r *router) AddHandler(prefix string, path *regexp.Regexp, fn http.HandlerFunc, methods ...string) error {
	// If methods is empty, default to GET
	if len(methods) == 0 {
		methods = []string{"GET"}
	}

	// Append the route
	r.routes = append(r.routes, route{normalizePath(prefix, true), path, fn, methods})

	// Sort routes by prefix length, longest first, and then by path != nil vs nil
	sort.Slice(r.routes, func(i, j int) bool {
		if len(r.routes[i].prefix) < len(r.routes[j].prefix) {
			return false
		}
		if len(r.routes[i].prefix) == len(r.routes[j].prefix) && r.routes[i].path == nil {
			return false
		}
		return true
	})

	// Return success
	return nil
}

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

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, params := r.get(req.Method, req.URL.Path)
	if route == nil {
		util.ServeError(w, http.StatusNotFound)
		return
	}

	// Check methods
	// TODO: This is not efficient
	for _, method := range route.methods {
		if req.Method == method {
			route.fn(w, req.Clone(context.WithPrefixParams(req.Context(), route.prefix, params)))
			return
		}
	}
	// Return method not allowed
	util.ServeError(w, http.StatusMethodNotAllowed)
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// get returns the route for the given path and method, and the parameters matched
// or returns nil for the route otherwise
func (r *router) get(method, path string) (*route, []string) {
	// Check cache
	if route, params := r.getcached(method, path); route != nil {
		return route, params
	}

	// Search routes
	for i := range r.routes {
		route := &r.routes[i]

		// Check against the prefix
		if !strings.HasPrefix(path, route.prefix) {
			continue
		}

		// Add a / to the beginning of the path
		relpath := normalizePath(path[len(route.prefix):], false)

		// Check for default route: this is the route that matches everything
		if route.path == nil {
			// Set cache
			r.setcached(method, path, i, nil)

			// Return route and params
			return route, nil
		}

		// Check for route with a regular expression
		if params := route.path.FindStringSubmatch(relpath); params != nil {
			// Set cache
			r.setcached(method, path, i, params[1:])

			// Return route and params
			return route, params[1:]
		}
	}

	// No match
	return nil, nil
}

// getcached returns the route for the given path, and the parameters matched
// or returns nil for the route otherwise
func (r *router) getcached(method, path string) (*route, []string) {
	r.RLock()
	defer r.RUnlock()
	cached, exists := r.cache[method+path]
	if !exists {
		return nil, nil
	} else {
		return &r.routes[cached.index], cached.matched
	}
}

// setcached puts a route into the cache
func (r *router) setcached(method, path string, index int, params []string) {
	r.Lock()
	defer r.Unlock()
	r.cache[method+path] = &cached{index, params}
}

// Add a / to the beginning and end of the path
func normalizePath(path string, end bool) string {
	if !strings.HasPrefix(path, pathSeparator) {
		path = pathSeparator + path
	}
	if end && !strings.HasSuffix(path, pathSeparator) {
		path = path + pathSeparator
	}
	return path
}
