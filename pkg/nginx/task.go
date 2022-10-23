package nginx

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run until done
func (r *nginx) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

// Return event channel
func (*nginx) C() <-chan Event {
	return nil
}
