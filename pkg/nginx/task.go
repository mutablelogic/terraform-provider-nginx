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
// TYPES

type configs struct {
	*nginx
	c map[string]plugin.NginxConfig
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run until done
func (r *nginx) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	configs := r.NewConfig()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := configs.run(); err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	}
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

func (r *nginx) NewConfig() *configs {
	return &configs{
		nginx: r,
		c:     map[string]plugin.NginxConfig{},
	}
}

func (r *configs) run() error {
	// TODO: Only enumerate on first run or if any folder has changed
	configs, err := r.Enumerate()
	if err != nil {
		return err
	}
	for _, config := range configs {
		name := config.Name()
		if _, exists := r.c[name]; !exists {
			fmt.Println("New config:", name)
			r.c[name] = config
		}
	}

	// Return success
	return nil
}
