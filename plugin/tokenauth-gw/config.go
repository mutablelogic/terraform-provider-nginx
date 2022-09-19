package main

import (
	"context"

	// Module imports
	plugin "github.com/mutablelogic/terraform-provider-nginx/plugin"
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/plugin/tokenauth"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label     string `json:"label,omitempty"`
	AuthLabel string `json:"auth,omitempty"`
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultLabel     = tokenauth.DefaultLabel + "-gw"
	defaultAuthLabel = tokenauth.DefaultLabel
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Set confuguration defaults
	if c.Label == "" {
		c.Label = defaultLabel
	}
	if c.AuthLabel == "" {
		c.AuthLabel = defaultAuthLabel
	}

	// Get token auth
	if task, ok := provider.TaskWithLabel(ctx, c.AuthLabel).(plugin.TokenAuth); !ok || task == nil {
		return nil, ErrNotFound.Withf("Not found: %q", c.AuthLabel)
	} else {
		return NewWithConfig(c, task)
	}
}

func (c Config) Name() string {
	return defaultLabel
}
