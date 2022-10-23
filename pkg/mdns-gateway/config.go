package mdns_gateway

import (
	"context"
	"time"

	// Modules
	mdns "github.com/mutablelogic/terraform-provider-nginx/pkg/mdns"
	discovery "github.com/mutablelogic/terraform-provider-nginx/pkg/mdns-discovery"
	types "github.com/mutablelogic/terraform-provider-nginx/pkg/types"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	L         string         `hcl:"label,label" json:"label,omitempty"`            // Label
	D         string         `hcl:"domain,optional" json:"domain,omitempty"`       // Domain (defaults to local)
	Interface string         `hcl:"interface,optional" json:"interface,omitempty"` // Interface name
	T         types.Duration `hcl:"ttl,optional" json:"ttl,omitempty"`             // TTL
	Task      types.Task     `hcl:"discovery,optional" json:"discovery,omitempty"` // NetServicesTask (optional)
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultName    = mdns.DefaultName + "-gateway"
	DefaultTimeout = 5 * time.Second
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) Name() string {
	return DefaultName
}

func (c Config) Label() string {
	if c.L == "" {
		return mdns.DefaultLabel
	} else {
		return c.L
	}
}

func (c Config) Domain() string {
	if c.D == "" {
		return mdns.DefaultLabel
	} else {
		return c.D
	}
}

func (c Config) TTL() time.Duration {
	if time.Duration(c.T) <= 0 {
		return discovery.DefaultTTL
	} else {
		return time.Duration(c.T)
	}
}

// Return a new task. Label for the task can be retrieved from context
func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Check label and domain
	if !util.IsIdentifier(c.Label()) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label())
	}
	if !util.IsIdentifier(c.Domain()) {
		return nil, ErrBadParameter.Withf("domain: %q", c.Domain())
	}

	// Create NetServiceTask
	if c.Task.Task == nil {
		if mdns, err := provider.New(ctx, mdns.Config{
			L:         c.Label(),
			Interface: c.Interface,
		}); err != nil {
			return nil, err
		} else if task, err := provider.New(ctx, discovery.Config{
			L:   c.Label(),
			D:   c.Domain(),
			TTL: types.Duration(c.TTL()),
			T:   types.Task{Task: mdns},
		}); err != nil {
			return nil, err
		} else {
			c.Task.Task = task
		}
	}

	// Create NetDiscovery task
	if c.Task.Task == nil {
		if task, err := provider.New(ctx, mdns.Config{
			L:         c.Label(),
			Interface: c.Interface,
		}); err != nil {
			return nil, err
		} else {
			c.Task.Task = task
		}
	}

	// Check task
	if _, ok := c.Task.Task.(NetServiceTask); !ok {
		return nil, ErrBadParameter.Withf("mdns: %v", c.Task.Task)
	}

	// Return configuration
	return NewWithConfig(c)
}
