package main

import (
	"context"

	// Module imports
	tokenauth "github.com/mutablelogic/terraform-provider-nginx/plugin/tokenauth"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label  string `json:"label,omitempty"`
	Auth   Task   `json:"-"` // plugin.TokenAuth
	Router Task   `json:"-"` // plugin.Router
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultLabel = tokenauth.DefaultLabel + "-gw"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Set confuguration defaults
	if c.Label == "" {
		c.Label = defaultLabel
	}

	// Return new task
	return NewWithConfig(c)
}

func (c Config) Name() string {
	return defaultLabel
}
