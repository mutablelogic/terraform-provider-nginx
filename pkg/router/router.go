package router

import (
	"net/http"
	"regexp"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type router struct {
	http.Handler
}

/*
   API:

   GET /           Returns the list of available configurations
   GET /:name      Returns a specific configuration
   POST /:name	   Creates a new configuration or updates an existing one, requires the configuraiton in the body.
                   This does not enable the configuration.
   DELETE /:name   Removes a configuration
   PATCH /:name    Enables or diables a configuration { enabled: true } or { enabled: false } as body
*/

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	reList   = regexp.MustCompile(`^/$`)
	reObject = regexp.MustCompile(`^/([a-zA-Z0-9_\-]+)$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewRouter() *router {
	return &router{}
}

/////////////////////////////////////////////////////////////////////
// MATCHER
