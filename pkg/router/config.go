package router

import (
	"context"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	L string `hcl:"label,label" json:"label,omitempty"`
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabel  = "router"
	pathSeparator = "/"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) Name() string {
	return DefaultLabel
}

func (c Config) Label() string {
	if c.L == "" {
		return DefaultLabel
	} else {
		return c.L
	}
}

// Return a new task. Label for the task can be retrieved from context
func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	if c.L == "" {
		c.L = DefaultLabel
	}
	if !util.IsIdentifier(c.Label()) {
		return nil, ErrBadParameter.Withf("label: %q", c.L)
	}

	// Return configuration
	return NewWithConfig(c)
}
