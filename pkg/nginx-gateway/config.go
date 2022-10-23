package nginx_gateway

import (
	"context"

	// Module imports

	"github.com/mutablelogic/terraform-provider-nginx/pkg/nginx"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/types"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label_ string     `json:"label,omitempty"`
	Prefix string     `json:"prefix,omitempty"`
	Nginx  types.Task `json:"nginx"`  // plugin.Nginx
	Router types.Task `json:"router"` // plugin.Router
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
	if _, ok := c.Router.Task.(Router); c.Router.Task == nil || !ok {
		return nil, ErrBadParameter.With("router")
	}
	if _, ok := c.Nginx.Task.(Nginx); c.Nginx.Task == nil || !ok {
		return nil, ErrBadParameter.With("nginx")
	}

	// Set configuration defaults
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

func (c Config) Label() string {
	if c.Label_ == "" {
		return DefaultLabel
	} else {
		return c.Label_
	}
}
