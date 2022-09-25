package nginx_gateway

import (
	"context"

	// Module imports

	"github.com/mutablelogic/terraform-provider-nginx/pkg/nginx"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label  string `json:"label,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Nginx  Task   `json:"-"` // plugin.Nginx
	Router Task   `json:"-"` // plugin.Router
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabelSuffix = "-gw"
	DefaultPathSuffix  = "/v1"
	DefaultLabel       = nginx.DefaultLabel + DefaultLabelSuffix
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Check arguments
	if _, ok := c.Router.(Router); c.Router == nil || !ok {
		return nil, ErrBadParameter.With("router")
	}
	if _, ok := c.Nginx.(TokenAuth); c.Nginx == nil || !ok {
		return nil, ErrBadParameter.With("nginx")
	}

	// Set confuguration defaults
	if c.Label == "" {
		c.Label = c.Nginx.Label() + DefaultLabelSuffix
	}
	if c.Prefix == "" {
		c.Prefix = "/" + c.Nginx.Label() + DefaultPathSuffix
	}

	// Check parameters
	if !util.IsIdentifier(c.Label) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label)
	}

	// Return new task
	return NewWithConfig(c)
}

func (c Config) Name() string {
	return DefaultLabel
}
