package tokenauth_gateway

import (
	"fmt"
	"net/http"
	"regexp"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"

	// Module imports
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type gateway struct {
	TokenAuth
	label, prefix string
	middleware    []string
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	rePathList         = regexp.MustCompile(`^/$`)
	rePathCreateRevoke = regexp.MustCompile(`^/(` + util.ReIdentifier + `)/?$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (Task, error) {
	plugin := new(gateway)
	plugin.label = c.Label
	plugin.prefix = c.Prefix
	plugin.TokenAuth = c.Auth.(TokenAuth)

	// Register middleware
	if err := c.Router.(Router).AddMiddleware(MiddlewareName, plugin.AuthenticateHandler); err != nil {
		return nil, err
	}
	if err := c.Router.(Router).AddMiddleware(MiddlewareAdminName, plugin.AuthenticateAdminHandler); err != nil {
		return nil, err
	}

	// Register handlers
	if err := c.Router.(Router).AddHandler(c.Prefix, rePathList, plugin.ListHandler, http.MethodGet); err != nil {
		return nil, err
	}
	if err := c.Router.(Router).AddHandler(c.Prefix, rePathCreateRevoke, plugin.CreateHandler, http.MethodPost); err != nil {
		return nil, err
	}
	if err := c.Router.(Router).AddHandler(c.Prefix, rePathCreateRevoke, plugin.RevokeHandler, http.MethodDelete); err != nil {
		return nil, err
	}

	// Return success
	return plugin, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (plugin *gateway) String() string {
	str := "<tokenauth-gateway"
	str += fmt.Sprintf(" label=%q", plugin.label)
	str += fmt.Sprintf(" prefix=%q", plugin.prefix)
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (plugin *gateway) Prefix() string {
	return plugin.prefix
}

func (plugin *gateway) Middleware() []string {
	return plugin.middleware
}
