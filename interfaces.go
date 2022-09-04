package nginxgw

import (
	"context"
	"net/http"
	"regexp"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Provider runs many tasks simultaneously
type Provider interface {
	Task

	// Create a new task
	New(ctx context.Context, config TaskPlugin) (Task, error)

	// Return a task with the given label
	TaskWithLabel(context.Context, string) Task

	// Return tasks with the given name
	TasksWithName(context.Context, string) []Task
}

// TaskPlugin provides methods to register a Task
type TaskPlugin interface {
	// Return the name of the task
	Name() string

	// Return a new task. Label for the task can be retrieved from context
	New(context.Context, Provider) (Task, error)
}

// Task runs a single task, whilst emitting events
type Task interface {
	// Return unique label for the task
	Label() string

	// Run is called to start the task and block until context is cancelled
	Run(context.Context) error

	// C returns a channel on which events can be received, or returns nil
	// if the task does not emit events
	C() <-chan Event
}

// Event will emit key/value pairs or errors emited on a channel
type Event interface {
	Key() any
	Value() any
	Error() error

	// Event can be emitted to a channel. Returns false if unable to do so
	Emit(chan<- Event) bool
}

// Router is a task which maps paths to routes
type Router interface {
	// Add a prefix/path mapping to a handler for one or more HTTP methods
	AddHandler(prefix string, path *regexp.Regexp, fn http.HandlerFunc, methods ...string) error

	// Add middleware for a unique name
	AddMiddleware(name string, fn func(http.HandlerFunc) http.HandlerFunc) error

	// Set middleware for a prefix. Called from left to right.
	SetMiddleware(prefix string, chain ...string) error
}
