package nginx_gateway

import (
	"context"
	"fmt"
	"time"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/event"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Run until done
func (plugin *gateway) Run(ctx context.Context) error {
	ticker := time.NewTimer(100 * time.Millisecond)
	defer ticker.Stop()
	defer plugin.stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := plugin.enumerate(); err != nil {
				event.NewError(err).Emit(plugin.ch)
			}
			ticker.Reset(time.Second)
		}
	}
}

// Return label
func (plugin *gateway) Label() string {
	return plugin.label
}

// Return event channel
func (plugin *gateway) C() <-chan Event {
	return plugin.ch
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (plugin *gateway) stop() {
	close(plugin.ch)
}

func (plugin *gateway) enumerate() error {
	// TODO: Only enumerate on first run or if any folder has changed
	configs, err := plugin.Enumerate()
	if err != nil {
		return err
	}
	for _, config := range configs {
		fmt.Println("emit", config)
		name := config.Name()
		event.NewEvent(name, config).Emit(plugin.ch)
	}
	return nil
}
