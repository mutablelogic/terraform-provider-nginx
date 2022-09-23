package router

import (
	"context"
	"regexp"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label string `hcl:"label,label"`
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabel = "router"
)

var (
	reValidName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]+$`)
)

const (
	pathSeparator = "/"
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
	if reValidName.MatchString(c.Label) == false {
		return nil, ErrBadParameter.Withf("label: %q", c.Label)
	}

	// Return configuration
	return NewWithConfig(c)
}
