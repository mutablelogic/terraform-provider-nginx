package logger

import (
	"context"

	// Modules
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label_ string `hcl:"label,label" json:"label,omitempty"` // Label for the configuration
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabel = "logger"
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
	// Check label
	if !util.IsIdentifier(c.Label()) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label())
	}

	// Return configuration
	return NewWithConfig(c)
}
