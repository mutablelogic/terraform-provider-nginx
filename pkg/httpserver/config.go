package httpserver

import (
	"context"
	"time"

	// Modules
	router "github.com/mutablelogic/terraform-provider-nginx/pkg/router"
	types "github.com/mutablelogic/terraform-provider-nginx/pkg/types"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label_  string         `hcl:"label,label" json:"label"`
	Router  types.Task     `hcl:"router,optional" json:"router"`
	Addr    string         `hcl:"listen,optional" json:"listen"`   // Address or path for binding HTTP server
	TLS     *TLS           `hcl:"tls,block" json:"tls"`            // TLS parameters
	Timeout types.Duration `hcl:"timeout,optional" json:"timeout"` // Read timeout on HTTP requests
}

type TLS struct {
	Key  string `hcl:"key"`  // Path to TLS Private Key
	Cert string `hcl:"cert"` // Path to TLS Certificate
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabel   = "httpserver"
	DefaultTimeout = 10 * time.Second
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

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

// Return a new task. Label for the task can be retrieved from context
func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Set timeout
	if c.Timeout == 0 {
		c.Timeout = types.Duration(DefaultTimeout)
	}
	// Check label
	if !util.IsIdentifier(c.Label()) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label())
	}
	// Create a router if it's not provided
	if c.Router.Task == nil {
		if router, err := provider.New(ctx, router.Config{
			Label_: c.Label() + "-router",
		}); err != nil {
			return nil, err
		} else {
			c.Router.Task = router
		}
	}

	// Return configuration
	return NewWithConfig(c)
}
