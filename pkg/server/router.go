package server

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type router struct {
	sync.RWMutex
	routes []route
	cache  map[string]*cached
}

type cached struct {
	index   int
	matched []string
}

type route struct {
	path    *regexp.Regexp
	fn      http.HandlerFunc
	methods []string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewRouter() *router {
	this := new(router)
	this.cache = make(map[string]*cached)
	return this
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (r *router) AddRoute(path *regexp.Regexp, fn http.HandlerFunc, methods ...string) {
	// If methods is empty, default to GET
	if len(methods) == 0 {
		methods = []string{"GET"}
	}
	r.routes = append(r.routes, route{path, fn, methods})
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, params := r.get(req.URL.Path)
	if route == nil {
		ServeError(w, http.StatusNotFound)
		return
	}
	// Check methods
	for _, method := range route.methods {
		if req.Method == method {
			route.fn(w, req.Clone(ctxWithParams(req.Context(), params)))
			return
		}
	}
	// Return method not allowed
	ServeError(w, http.StatusMethodNotAllowed)
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// get returns the route for the given path, and the parameters matched
// or returns nil for the route otherwise
func (r *router) get(path string) (*route, []string) {

	// Check cache
	if route, params := r.getcached(path); route != nil {
		return route, params
	}

	// Search routes
	for i := range r.routes {
		fmt.Println("Check", r.routes[i].path)
		route := &r.routes[i]
		if params := route.path.FindStringSubmatch(path); params != nil {
			// Set cache
			r.setcached(path, i, params[1:])
			// Return route and params
			return route, params[1:]
		}
	}

	// No match
	return nil, nil
}

// getcached returns the route for the given path, and the parameters matched
// or returns nil for the route otherwise
func (r *router) getcached(path string) (*route, []string) {
	r.RLock()
	defer r.RUnlock()

	cached, exists := r.cache[path]
	if !exists {
		return nil, nil
	} else {
		return &r.routes[cached.index], cached.matched
	}
}

// setcached puts a route into the cache
func (r *router) setcached(path string, index int, params []string) {
	r.Lock()
	defer r.Unlock()

	r.cache[path] = &cached{index, params}
}
