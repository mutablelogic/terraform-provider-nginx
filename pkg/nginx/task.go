package nginx

import (
	"context"
	"fmt"
	"time"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"

	// Modules
	"github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run until done
func (r *nginx) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	configs := map[string]plugin.NginxConfig{}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.run(configs); err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	}

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

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (r *nginx) run(configs map[string]plugin.NginxConfig) error {
	return nil
}
