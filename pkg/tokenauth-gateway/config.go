package tokenauth_gateway

import (
	"context"

	// Module imports

	"github.com/mutablelogic/terraform-provider-nginx/pkg/tokenauth"
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
	Auth   Task   `json:"-"` // plugin.TokenAuth
	Router Task   `json:"-"` // plugin.Router
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabelSuffix = "-gw"
	DefaultPathSuffix  = "/v1"
	DefaultLabel       = tokenauth.DefaultLabel + DefaultLabelSuffix
	MiddlewareName     = tokenauth.DefaultLabel
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Check arguments
	if _, ok := c.Router.(Router); c.Router == nil || !ok {
		return nil, ErrBadParameter.With("router")
	}
	if _, ok := c.Auth.(TokenAuth); c.Auth == nil || !ok {
		return nil, ErrBadParameter.With("auth")
	}

	// Set confuguration defaults
	if c.Label == "" {
		c.Label = c.Auth.Label() + DefaultLabelSuffix
	}
	if c.Prefix == "" {
		c.Prefix = "/" + c.Auth.Label() + DefaultPathSuffix
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
