package tokenauth

import (
	"context"
	"os"
	"path/filepath"
	"time"

	// Module imports
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label string        `json:"label,omitempty"`
	Path  string        `json:"path,omitempty"`
	File  string        `json:"file,omitempty"`
	Delta time.Duration `json:"delta,omitempty"`
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabel = "tokenauth"
	AdminToken   = "admin"
)

const (
	defaultFile                 = DefaultLabel + ".json"
	defaultLength               = 32
	defaultDelta                = time.Second * 30
	defaultEventChannelCapacity = 1000
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Set confuguration defaults
	if c.Label == "" {
		c.Label = DefaultLabel
	}
	if c.File == "" {
		c.File = defaultFile
	}
	if c.Delta <= 0 {
		c.Delta = defaultDelta
	}

	// Check label is valid
	if !util.IsIdentifier(c.Label) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label)
	}

	// If path is empty, then use the default and maybe create it
	if c.Path == "" {
		if path, err := os.UserConfigDir(); err != nil {
			return nil, err
		} else {
			c.Path = filepath.Join(path, c.Label)
		}
		if err := os.MkdirAll(c.Path, 0755); err != nil {
			return nil, err
		}
	}

	// Return configuration
	return NewWithConfig(c)
}

func (c Config) Name() string {
	return DefaultLabel
}
