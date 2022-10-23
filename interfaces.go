package nginxgw

import (
	"context"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Provider runs many tasks simultaneously
type Provider interface {
	Task

	// Create a new task
	New(ctx context.Context, config TaskPlugin) (Task, error)
}

// TaskPlugin provides methods to register a Task
type TaskPlugin interface {
	Name() string  // Return the name of the task
	Label() string // Return the label of the task

	// Return a new task. Label for the task can be retrieved from context
	New(context.Context, Provider) (Task, error)
}

// Task runs a single task, whilst emitting events
type Task interface {
	// Run is called to start the task and block until context is cancelled
	Run(context.Context) error

	// Sub returns a channel on which events can be received, or returns nil
	// if the task does not emit events
	Sub() <-chan Event

	// Unsub is called to unsubscribe from an existing events channel
	Unsub(<-chan Event)
}

// Event will emit key/value pairs or errors emited on a channel
type Event interface {
	Key() any
	Value() any
	Error() error

	// Event can be emitted to a channel. Returns false if unable to do so
	Emit(chan<- Event) bool
}
