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
	Label string `hcl:"label,label"` // Label for the configuration
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

// Return a new task. Label for the task can be retrieved from context
func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Set label
	if c.Label == "" {
		c.Label = DefaultLabel
	}

	// Check label
	if !util.IsIdentifier(c.Label) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label)
	}

	// Return configuration
	return NewWithConfig(c)
}
