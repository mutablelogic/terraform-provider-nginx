package nginx

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run until done
func (*nginx) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// Return label
func (r *nginx) Label() string {
	return r.label
}

// Return event channel
func (*nginx) C() <-chan Event {
	return nil
}
