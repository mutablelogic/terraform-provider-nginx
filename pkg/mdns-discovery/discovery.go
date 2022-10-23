package mdns_discovery

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	// Modules
	event "github.com/mutablelogic/terraform-provider-nginx/pkg/event"
	mdns "github.com/mutablelogic/terraform-provider-nginx/pkg/mdns"
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	//. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type task struct {
	provider.Task
	mdns   mdns.DNSTask
	domain string
	ttl    time.Duration
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (NetServiceTask, error) {
	t := new(task)
	t.mdns = c.T.Task.(mdns.DNSTask)
	t.domain = util.Fqn(c.Domain())
	t.ttl = time.Duration(c.TTL)

	// Return success
	return t, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t *task) String() string {
	str := "<mdns.discovery"
	if t.domain != "" {
		str += fmt.Sprintf(" domain=%q", t.domain)
	}
	if t.ttl != 0 {
		str += fmt.Sprintf(" ttl=%v", t.ttl)
	}
	if t.mdns != nil {
		str += fmt.Sprintf(" mdns=%v", t.mdns)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t *task) Run(ctx context.Context) error {
	// Receive events from mdns
	ch := t.mdns.Sub()
	defer t.mdns.Unsub(ch)
	for {
		select {
		case <-ctx.Done():
			// Close channels and end
			t.Emit(nil)
			return ctx.Err()
		case event := <-ch:
			// Parse events which are ReceiveIP4 and ReceiveIP6
			if event.Key() == Receive {
				t.parse(event.Value().(mdns.Message))
			}
		}
	}
}

func (t *task) Discover(ctx context.Context) ([]string, error) {
	var result = map[string]bool{}

	// Subscribe to receive discovery messages
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch := t.Sub()
	FOR_LOOP:
		for {
			select {
			case <-ctx.Done():
				t.Unsub(ch)
				break FOR_LOOP
			case event := <-ch:
				if service := t.serviceFromEvent(event); service != "" {
					result[service] = true
				}
			}
		}
	}()

	// Create a message to send - without defined interface
	if message, err := mdns.MessageWithQuestion(util.Fqn(ServicesQuery, t.domain), net.Interface{}); err != nil {
		return nil, err
	} else if err := t.mdns.Send(ctx, message); err != nil {
		return nil, err
	}

	// Wait for goroutine to end
	wg.Wait()

	// Return success
	return arrayOfKeys(result), nil
}

// Resolve service name into service instances
func (t *task) Resolve(ctx context.Context, service string) ([]NetService, error) {
	// Create a message to send
	if message, err := mdns.MessageWithQuestion(util.Fqn(service, t.domain), net.Interface{}); err != nil {
		return nil, err
	} else if err := t.mdns.Send(ctx, message); err != nil {
		return nil, err
	}

	// Return success
	return nil, nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (t *task) parse(message mdns.Message) {
	for _, ptr := range message.PTR() {
		// Check domain
		service := ptr.Service()
		if !strings.HasSuffix(service, t.domain) {
			return
		} else {
			service = util.Unfqn(service, t.domain)
		}
		// Service record
		if service == ServicesQuery && ptr.TTL() > 0 {
			t.Emit(event.NewEvent(Discovered, ptr.Name()))
		} else {
			fmt.Println("TODO:", ptr)
		}
	}
}

func (t *task) serviceFromEvent(event Event) string {
	if event == nil || event.Key() != Discovered {
		return ""
	} else if name, ok := event.Value().(string); !ok {
		return ""
	} else if !strings.HasSuffix(name, t.domain) {
		return ""
	} else {
		return util.Unfqn(name, t.domain)
	}
}

func arrayOfKeys(in map[string]bool) []string {
	result := make([]string, 0, len(in))
	for key := range in {
		result = append(result, key)
	}
	return result
}
