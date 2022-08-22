package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	// Module imports
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type runner struct {
	*config
	available map[string]*record
	c         chan Event
}

type record struct {
	file    File
	hash    string
	enabled bool
	marked  bool
}

type EventType uint

type Event struct {
	Type  EventType // Type associated with the event
	File  File      // File associated with the event
	Error error     // Any errors
}

type Object struct {
	Name    string `json:"name"`
	Path    string `json:"path,omitempty"`
	Enabled bool   `json:"enabled"`
	Body    []byte `json:"body,omitempty"`
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EventNone EventType = iota
	EventError
	EventCreate
	EventRemove
	EventChange
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (cfg Config) NewRunner() (*runner, error) {
	config, err := cfg.New()
	if err != nil {
		return nil, err
	}
	return &runner{
		config:    config,
		available: make(map[string]*record),
		c:         make(chan Event, 1000),
	}, nil
}

func (r *runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			changed, err := r.Changed()
			if err != nil {
				r.emit(Event{Type: EventError, Error: err})
			} else if changed {
				if err := r.Enumerate(); err != nil {
					r.emit(Event{Type: EventError, Error: err})
				}
			}
		case <-ctx.Done():
			close(r.c)
			if err := ctx.Err(); err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			} else {
				return err
			}
		}
	}
}

func (r *runner) Enumerate() error {
	// Enumerate all existing files
	files, err := r.EnumerateAvailable()
	if err != nil {
		return err
	}

	// Mark existing records so we can remove them later
	for _, record := range r.available {
		record.marked = true
	}

	var names = make(map[string]bool)
	var result error
	for _, f := range files {
		key := f.Name()
		// Report on any duplicate names
		if _, exists := names[key]; exists {
			result = multierror.Append(result, ErrDuplicateEntry.With(key))
			continue
		} else {
			names[key] = true
		}
		// If no record exists for this file, create one
		if _, exists := r.available[key]; !exists {
			r.available[key] = &record{file: f}
			r.emit(Event{Type: EventCreate, File: f})
		}
		// Remove mark
		r.available[key].marked = false
	}

	// If there are any records which are marked, then remove them
	for _, record := range r.available {
		if record.marked {
			r.emit(Event{Type: EventRemove, File: record.file})
			delete(r.available, record.file.Name())
		}
	}

	// Return success
	return nil
}

// List returns all configurations
func (r *runner) List() []Object {
	result := []Object{}
	for _, record := range r.available {
		result = append(result, Object{
			Name:    record.file.Name(),
			Enabled: record.enabled,
		})
	}
	return result
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v EventType) String() string {
	switch v {
	case EventNone:
		return "EventNone"
	case EventError:
		return "EventError"
	case EventCreate:
		return "EventCreate"
	case EventRemove:
		return "EventRemove"
	case EventChange:
		return "EventChange"
	default:
		return "???"
	}
}

func (e Event) String() string {
	str := "<event"
	if e.Type != EventNone {
		str += fmt.Sprint(" type=", e.Type)
	}
	if e.File != nil {
		str += fmt.Sprint(" file=", e.File)
	}
	if e.Error != nil {
		str += fmt.Sprint(" error=", e.Error)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (r *runner) C() <-chan Event {
	return r.c
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (r *runner) emit(e Event) {
	select {
	case r.c <- e:
		return
	default:
		return
	}
}
