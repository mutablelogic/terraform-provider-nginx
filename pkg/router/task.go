package router

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run until done
func (*router) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// Return label
func (r *router) Label() string {
	return r.label
}

// Return event channel. No events are sent by the router
func (*router) C() <-chan Event {
	return nil
}
