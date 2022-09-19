package main

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run will write the authorization tokens back to disk if they have been modified
func (plugin *gateway) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (plugin *gateway) Label() string {
	return c.label
}

func (plugin *gateway) C() <-chan Event {
	return nil
}
