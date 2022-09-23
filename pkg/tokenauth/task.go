package tokenauth

import (
	"context"
	"time"

	// Module imports
	event "github.com/mutablelogic/terraform-provider-nginx/pkg/event"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run will write the authorization tokens back to disk if they have been modified
func (c *auth) Run(ctx context.Context) error {
	ticker := time.NewTicker(c.delta)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(c.ch)
			_, err := c.writeIfModified()
			return err
		case <-ticker.C:
			c.Lock()
			if written, err := c.writeIfModified(); err != nil {
				event.NewError(err).Emit(c.ch)
			} else if written {
				event.NewEvent(nil, "Written tokens to disk").Emit(c.ch)
			}
			c.Unlock()
		}
	}
}

// Return label
func (c *auth) Label() string {
	return c.label
}

// Return event channel
func (c *auth) C() <-chan Event {
	return c.ch
}
