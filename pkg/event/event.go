package event

import "fmt"

type Event struct {
	Type  any   // Type associated with the event
	Value any   // Value associated with the event
	Error error // Any errors
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEvent(t any, v any) *Event {
	return &Event{Type: t, Value: v}
}

func NewError(err error) *Event {
	return &Event{Error: err}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e *Event) String() string {
	str := "<event"
	if e.Type != nil {
		str += fmt.Sprint(" type=", e.Type)
	}
	if e.Value != nil {
		switch e.Value.(type) {
		case string:
			str += fmt.Sprintf(" value=%q", e.Value)
		default:
			str += fmt.Sprint(" value=", e.Value)
		}
	}
	if e.Error != nil {
		str += fmt.Sprint(" error=", e.Error)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e *Event) Emit(ch chan<- *Event) bool {
	select {
	case ch <- e:
		return true
	default:
		return false
	}
}
