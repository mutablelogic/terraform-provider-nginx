package event

import (
	"fmt"
	"sync"

	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// PubSub is a publish/subscribe instance for events
type PubSub struct {
	sync.Mutex
	Cap uint
	ch  []chan Event
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

// Return string representation of a `PubSub` instance
func (p *PubSub) String() string {
	str := "<pubsub"
	if p.Cap > 0 {
		str += fmt.Sprint(" cap=", p.Cap)
	}
	if len(p.ch) > 0 {
		str += fmt.Sprint(" ch=", len(p.ch))
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// CREATE EVENT

// Sub creates a new subscriber channel and returns it
func (p *PubSub) Sub() <-chan Event {
	p.Lock()
	defer p.Unlock()
	for i := range p.ch {
		if p.ch[i] == nil {
			p.ch[i] = make(chan Event, p.Cap)
			return p.ch[i]
		}
	}
	p.ch = append(p.ch, make(chan Event, p.Cap))
	return p.ch[len(p.ch)-1]
}

// Unsub is called to unsubscribe a specific channel. Will panic
// if the channel is not subscribed
func (p *PubSub) Unsub(ch <-chan Event) {
	p.Lock()
	defer p.Unlock()
	for i := range p.ch {
		if ch == p.ch[i] {
			close(p.ch[i])
			p.ch[i] = nil
			return
		}
	}
	panic("Unsub called for unsubscribed channel")
}

// Emit can be called to send an event to all subscribers,
// and returns true if the event was sent to all channels
func (p *PubSub) Emit(e Event) bool {
	result := true
	for i, ch := range p.ch {
		if ch == nil {
			// do nothing if channel is closed
		} else if e == nil {
			close(ch)
			p.Lock()
			p.ch[i] = nil
			p.Unlock()
		} else if done := e.Emit(ch); !done {
			result = false
		}
	}
	if e == nil {
		p.Lock()
		p.ch = nil
		p.Unlock()
	}
	return result
}
