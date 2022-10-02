package provider

import (
	"context"
	"fmt"

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
	Label string `hcl:"label,label"`
}

// This is the most basic task, which just emits a single event on startup, and
// can be used as an example of task creation and lifecycle.
type Task struct {
	event.PubSub
	label string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Return the unique name of the task plugin
func (c Config) Name() string {
	return "task"
}

// Create a task instance from a configuration
func (c Config) New(context.Context, iface.Provider) (iface.Task, error) {
	t := new(Task)

	if !util.IsIdentifier(c.Label) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label)
	} else {
		t.label = c.Label
	}

	// Return success
	return t, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (task *Task) String() string {
	str := "<task"
	str += fmt.Sprintf(" label=%q", task.label)
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (task *Task) Label() string {
	return task.label
}

func (task *Task) SetLabel(value string) error {
	if !util.IsIdentifier(value) {
		return ErrBadParameter.Withf("label: %q", value)
	} else {
		task.label = value
	}

	// Return success
	return nil
}

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
