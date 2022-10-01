package event

import (
	"fmt"
	"strconv"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	key, value any
	err        error // Any errors
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEvent(key, value any) *event {
	return &event{key, value, nil}
}

func NewError(err error) *event {
	return &event{nil, nil, err}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e *event) String() string {
	str := "<event"
	if e.key != nil {
		str += fmt.Sprint(" key=", toString(e.key))
	}
	if e.value != nil {
		str += fmt.Sprint(" value=", toString(e.value))
	}
	if e.err != nil {
		str += fmt.Sprint(" error=", e.err)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (e *event) Key() any {
	return e.key
}

func (e *event) Value() any {
	return e.value
}

func (e *event) Error() error {
	return e.err
}

func (e *event) Emit(ch chan<- Event) bool {
	// Unbuffered channels block
	if cap(ch) == 0 {
		ch <- e
		return true
	}
	// Buffered channels don't block
	select {
	case ch <- e:
		return true
	default:
		return false
	}
}

/////////////////////////////////////////////////////////////////////
//PRIVATE METHODS

func toString(v any) string {
	switch v := v.(type) {
	case string:
		return strconv.Quote(v)
	default:
		return fmt.Sprint(v)
	}
}
