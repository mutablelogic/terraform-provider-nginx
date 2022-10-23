package mdns

import (
	"context"
	"net"
	"time"

	// Modules
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	L         string `hcl:"label,label" json:"label,omitempty"`            // Label for the configuration
	Interface string `hcl:"interface,optional" json:"interface,omitempty"` // The interface to bind to. Defaults to all broadcast interfaces

	iface []net.Interface
}

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultName  = "mdns"
	DefaultLabel = "local"
)

var (
	multicastAddrIp4 = &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: 5353}
	multicastAddrIp6 = &net.UDPAddr{IP: net.ParseIP("ff02::fb"), Port: 5353}
	sendRetryCount   = 3
	sendRetryDelta   = 100 * time.Millisecond
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Return the name of the plugin
func (c Config) Name() string {
	return DefaultName
}

// Return the label of the instance
func (c Config) Label() string {
	if c.L == "" {
		return DefaultLabel
	} else {
		return c.L
	}
}

// Return a new task. Label for the task can be retrieved from context
func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Check label
	if !util.IsIdentifier(c.Label()) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label())
	}

	// Check interfaces
	var i net.Interface
	var err error
	if c.Interface != "" {
		i, err = interfaceForName(c.Interface)
		if err != nil {
			return nil, err
		}
	}
	if ifaces, err := multicastInterfaces(i); err != nil {
		return nil, err
	} else if len(ifaces) == 0 {
		return nil, ErrBadParameter.With("no multicast interfaces defined")
	} else {
		c.iface = ifaces
	}

	// Return configuration
	return NewWithConfig(c)
}
