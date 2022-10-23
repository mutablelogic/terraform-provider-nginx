package mdns_discovery

import (
	"context"
	"time"

	// Modules
	mdns "github.com/mutablelogic/terraform-provider-nginx/pkg/mdns"
	types "github.com/mutablelogic/terraform-provider-nginx/pkg/types"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	L   string         `hcl:"label,label" json:"label,omitempty"`      // Label
	D   string         `hcl:"domain,optional" json:"domain,omitempty"` // mDNS Domain
	TTL types.Duration `hcl:"ttl,optional" json:"ttl,omitempty"`       // TTL
	T   types.Task     `hcl:"mdns,optional" json:"mdns,omitempty"`     // mDNS task
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultName   = mdns.DefaultName + "-discovery"
	ServicesQuery = "_services._dns-sd._udp"
	DefaultTTL    = 2 * time.Minute
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

// Return a new task. Label for the task can be retrieved from context
func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Check label and domain
	if !util.IsIdentifier(c.Label()) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label())
	}
	if !util.IsIdentifier(c.Domain()) {
		return nil, ErrBadParameter.Withf("domain: %q", c.Domain())
	}

	// Check ttl
	if c.TTL < 0 {
		return nil, ErrBadParameter.Withf("ttl: %v", c.TTL)
	} else if c.TTL == 0 {
		c.TTL = types.Duration(DefaultTTL)
	}

	// Create task
	if c.T.Task == nil {
		if task, err := provider.New(ctx, mdns.Config{
			L: c.Label(),
		}); err != nil {
			return nil, err
		} else {
			c.T.Task = task
		}
	}

	// Check task
	if _, ok := c.T.Task.(mdns.DNSTask); !ok {
		return nil, ErrBadParameter.Withf("mdns: %v", c.T.Task)
	}

	// Return configuration
	return NewWithConfig(c)
}
