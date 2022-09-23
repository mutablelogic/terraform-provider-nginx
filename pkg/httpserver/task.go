package httpserver

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run until done
func (*httpserver) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// Return label
func (r *httpserver) Label() string {
	return r.label
}

// Return event channel. No events are sent by the httpserver
func (*httpserver) C() <-chan Event {
	return nil
}
