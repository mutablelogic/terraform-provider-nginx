package nginx_gateway

import (
	"fmt"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
	// Module imports
)

/////////////////////////////////////////////////////////////////////
// TYPES

type gateway struct {
	Nginx
	label, prefix string
	middleware    []string
	ch            chan Event
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
// rePathList         = regexp.MustCompile(`^/$`)
// rePathCreateRevoke = regexp.MustCompile(`^/(` + util.ReIdentifier + `)/?$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (Task, error) {
	plugin := new(gateway)
	plugin.label = c.Label
	plugin.prefix = c.Prefix
	plugin.ch = make(chan Event, 100)
	plugin.Nginx = c.Nginx.(Nginx)

	// Register handlers
	//if err := c.Router.(Router).AddHandler(c.Prefix, rePathList, plugin.ListHandler, http.MethodGet); err != nil {
	//	return nil, err
	//}

	// Return success
	return plugin, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (plugin *gateway) String() string {
	str := "<nginx-gateway"
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
