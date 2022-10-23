package mdns

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	// Modules
	multierror "github.com/hashicorp/go-multierror"
	event "github.com/mutablelogic/terraform-provider-nginx/pkg/event"
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
	ipv4 "golang.org/x/net/ipv4"
	ipv6 "golang.org/x/net/ipv6"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

///////////////////////////////////////////////////////////////////////////////
// MESSAGE INTERFACE

// A `DNSTask` which receives and sends messages
type DNSTask interface {
	Task

	// Send a DNS message
	Send(context.Context, Message) error
}

///////////////////////////////////////////////////////////////////////////////
// TYPES

type mdns struct {
	provider.Task

	// Interface parameters
	ifaces map[int]net.Interface

	// Send parameters
	ch    chan Message
	count int
	delta time.Duration
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a task instance from a `Config` object
func NewWithConfig(c Config) (DNSTask, error) {
	t := new(mdns)
	t.ch = make(chan Message)
	t.count = sendRetryCount
	t.delta = sendRetryDelta

	// Make index->interface map
	t.ifaces = make(map[int]net.Interface, len(c.iface))
	for _, iface := range c.iface {
		t.ifaces[iface.Index] = iface
	}

	// Return success
	return t, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t *mdns) String() string {
	str := "<dns.task"
	if len(t.ifaces) > 0 {
		var ifaces []string
		for _, iface := range t.ifaces {
			ifaces = append(ifaces, iface.Name)
		}
		str += fmt.Sprintf(" ifaces=%q", ifaces)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t *mdns) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	var result error
	var ip4 *ipv4.PacketConn
	var ip6 *ipv6.PacketConn

	// Bind interfaces to sockets
	var ifaces []net.Interface
	for _, iface := range t.ifaces {
		ifaces = append(ifaces, iface)
	}

	// IP4
	if ip4, result = bindUdp4(ifaces, multicastAddrIp4); result != nil {
		return result
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.run4(ctx, ip4)
		}()
	}

	// IP6
	if ip6, result = bindUdp6(ifaces, multicastAddrIp6); result != nil {
		return result
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.run6(ctx, ip6)
		}()
	}

	// Send messages until done
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			close(t.ch)
			break FOR_LOOP
		case msg := <-t.ch:
			if msg != nil {
				if err := t.send(ip4, ip6, msg); err != nil {
					t.Emit(event.NewError(err))
				} else {
					t.Emit(event.NewEvent(Send, msg))
				}
			}
		}
	}

	// Close connections - these will quit tuen run4/run6 goroutines
	if err := ip4.Close(); err != nil {
		result = multierror.Append(result, err)
	}
	if err := ip6.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	// Wait until receive loops have completed
	wg.Wait()

	// Close subscriber channels
	t.Emit(nil)

	// Return any errors
	return result
}

func (t *mdns) Send(ctx context.Context, msg Message) error {
	timer := time.NewTimer(1 * time.Nanosecond)
	defer timer.Stop()

	i := 0
	for {
		i++
		select {
		case <-timer.C:
			t.ch <- msg
			if i > t.count {
				return nil
			}
			timer.Reset(t.delta * time.Duration(i))
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (t *mdns) interfaceWithIndex(ifIndex int) net.Interface {
	if iface, exists := t.ifaces[ifIndex]; exists {
		return iface
	} else {
		return net.Interface{}
	}
}

func (t *mdns) run4(ctx context.Context, conn *ipv4.PacketConn) {
	buf := make([]byte, 65536)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if n, cm, from, err := conn.ReadFrom(buf); err != nil {
				continue
			} else if cm == nil {
				continue
			} else if msg, err := MessageFromPacket(buf[:n], from, t.interfaceWithIndex(cm.IfIndex)); err != nil {
				t.Emit(event.NewError(err))
			} else if msg.IsAnswer() && len(msg.PTR()) > 0 {
				t.Emit(event.NewEvent(Receive, msg))
			}
		}
	}
}

func (t *mdns) run6(ctx context.Context, conn *ipv6.PacketConn) {
	buf := make([]byte, 65536)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if n, cm, from, err := conn.ReadFrom(buf); err != nil {
				continue
			} else if cm == nil {
				continue
			} else if msg, err := MessageFromPacket(buf[:n], from, t.interfaceWithIndex(cm.IfIndex)); err != nil {
				t.Emit(event.NewError(err))
			} else if msg.IsAnswer() && len(msg.PTR()) > 0 {
				t.Emit(event.NewEvent(Receive, msg))
			}
		}
	}
}

// Send a single DNS message to a particular interface or all interfaces if 0
func (t *mdns) send(conn4 *ipv4.PacketConn, conn6 *ipv6.PacketConn, message Message) error {
	var result error

	// Pack the message
	ifIndex := message.IfIndex()
	data, err := message.Bytes()
	if err != nil {
		return err
	}

	// IP4 send
	if conn4 != nil {
		var cm ipv4.ControlMessage
		if ifIndex != 0 {
			cm.IfIndex = ifIndex
			if _, err := conn4.WriteTo(data, &cm, multicastAddrIp4); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			for _, intf := range t.ifaces {
				cm.IfIndex = intf.Index
				if _, err := conn4.WriteTo(data, &cm, multicastAddrIp4); err != nil {
					result = multierror.Append(result, err)
				}
			}
		}
	}

	// IP6 send
	if conn6 != nil {
		var cm ipv6.ControlMessage
		if ifIndex != 0 {
			cm.IfIndex = ifIndex
			if _, err := conn6.WriteTo(data, &cm, multicastAddrIp6); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			for _, intf := range t.ifaces {
				cm.IfIndex = intf.Index
				if _, err := conn6.WriteTo(data, &cm, multicastAddrIp6); err != nil {
					result = multierror.Append(result, err)
				}
			}
		}
	}

	// Return any errors
	return result
}
