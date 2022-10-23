package mdns_gateway

import (
	// Modules
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/event"
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"

	// Namespace imports
	//. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type task struct {
	provider.Task
	netservice NetServiceTask
	ttl        time.Duration
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (NetServiceGateway, error) {
	t := new(task)
	t.netservice = c.Task.Task.(NetServiceTask)
	t.ttl = c.TTL()

	// Return success
	return t, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t *task) String() string {
	str := "<mdns.gateway"
	if t.ttl != 0 {
		str += fmt.Sprintf(" ttl=%v", t.ttl)
	}
	if t.netservice != nil {
		str += fmt.Sprintf(" netservice=%v", t.netservice)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t *task) Run(ctx context.Context) error {
	// ticker will be used to enumerate services
	ticker := time.NewTimer(time.Second)
	defer ticker.Stop()

	// Receive events from mdns
	for {
		select {
		case <-ctx.Done():
			// Close channels and end
			t.Emit(nil)
			return ctx.Err()
		case <-ticker.C:
			go func() {
				if err := t.discover(ctx); err != nil {
					t.Emit(event.NewError(err))
				}
				ticker.Reset(t.ttl)
			}()
		}
	}
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (plugin *task) Prefix() string {
	return "TODO"
}

func (plugin *task) Middleware() []string {
	return nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (t *task) discover(parent context.Context) error {
	var result error

	ctx, cancel := context.WithTimeout(parent, DefaultTimeout)
	defer cancel()
	services, err := t.netservice.Discover(ctx)
	if err != nil {
		return err
	}
	for _, service := range services {
		ctx, cancel := context.WithTimeout(parent, DefaultTimeout)
		defer cancel()
		if instances, err := t.netservice.Resolve(ctx, service); err != nil {
			result = multierror.Append(result, err)
		} else {
			for _, instance := range instances {
				t.Emit(event.NewEvent(Resolved, instance))
			}
		}
	}

	// Return success
	return result
}
