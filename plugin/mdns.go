package plugin

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// The message type
type MessageType uint

// Network Service
type NetService interface {
	Service() string
}

// Network Service Task
type NetServiceTask interface {
	Task

	// Discover service types on all interfaces
	Discover(context.Context) ([]string, error)

	// Resolve service name into service instances
	Resolve(context.Context, string) ([]NetService, error)
}

// Network Service Gateway
type NetServiceGateway interface {
	Gateway
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	Receive    MessageType = iota // A received message
	Send                          // A sent message
	Discovered                    // A service discovery message
	Resolved                      // A service was resolved into an instance
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v MessageType) String() string {
	switch v {
	case Receive:
		return "Receive"
	case Send:
		return "Send"
	case Discovered:
		return "Discovered"
	case Resolved:
		return "Resolved"
	default:
		return "[?? Invalid MessageType value]"
	}
}
