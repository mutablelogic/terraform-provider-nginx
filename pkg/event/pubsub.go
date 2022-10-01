package event

import (
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type PubSub struct {
	Cap uint
	ch  []chan Event
}

/////////////////////////////////////////////////////////////////////
// CREATE EVENT

func (p *PubSub) C() <-chan Event {
	// Create channel and return it
	p.ch = append(p.ch, make(chan Event, p.Cap))
	return p.ch[len(p.ch)-1]
}

func (p *PubSub) Emit(e Event) bool {
	result := true
	for _, ch := range p.ch {
		if e == nil {
			close(ch)
		} else if done := e.Emit(ch); !done {
			result = false
		}
	}
	if e == nil {
		p.ch = nil
	}
	return result
}
