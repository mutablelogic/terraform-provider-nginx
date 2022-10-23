package provider

import (
	"context"

	// Namespace imports
	. "github.com/djthorpe/go-errors"

	// Modules
	iface "github.com/mutablelogic/terraform-provider-nginx"
	event "github.com/mutablelogic/terraform-provider-nginx/pkg/event"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// Config is the basic configuration for a task, which includes a simple
// label for the task
type Config struct {
	Label_ string `hcl:"label,label" json:"label,omitempty"`
}

// This is the most basic task, which just emits a single event on startup, and
// can be used as an example of task creation and lifecycle.
type Task struct {
	event.PubSub
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultLabel = "main"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Return the unique name of the task plugin
func (c Config) Name() string {
	return "task"
}

// Return the label of the task instance
func (c Config) Label() string {
	if c.Label_ == "" {
		return defaultLabel
	} else {
		return c.Label_
	}
}

// Create a task instance from a configuration
func (c Config) New(context.Context, iface.Provider) (iface.Task, error) {
	t := new(Task)

	// Check label
	if !util.IsIdentifier(c.Label()) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label())
	}

	// Return success
	return t, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (task *Task) String() string {
	str := "<task"
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (task *Task) Run(ctx context.Context) error {
	// Emit starting event
	task.Emit(event.NewEvent("start", "started"))

	// Wait until done
	<-ctx.Done()

	// Close channels
	task.Emit(nil)

	// Return errors
	return ctx.Err()
}
