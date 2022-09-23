package tokenauth_gateway

import (
	"fmt"
	"net/http"
	"regexp"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type gateway struct {
	TokenAuth
	label, prefix string
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	rePathList = regexp.MustCompile(`^/$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (Task, error) {
	plugin := new(gateway)
	plugin.label = c.Label
	plugin.prefix = c.Prefix
	plugin.TokenAuth = c.Auth.(TokenAuth)

	// Register routes
	if err := c.Router.(Router).AddHandler(c.Prefix, rePathList, plugin.ListHandler, http.MethodGet); err != nil {
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
