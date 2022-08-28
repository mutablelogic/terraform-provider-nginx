package nginxgw

import (
	"context"
	"net/http"
	"regexp"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Kernel runs many tasks simultaneously
type Kernel interface {
	RouterTask

	// Run the kernel until cancel event
	Run(context.Context) error

	// Receive events
	C() <-chan Event

	// Add a Task to the kernel with a name
	Add(string, Task) error

	// Return a task with the given name
	Get(string) Task
}

// Task has a Run function which is called when the task is started, and it continues
// until context is cancelled. If can emit events to the channel.
type Task interface {
	Run(context.Context, Kernel) error
	C() <-chan Event
}

// Event can emit events to a channel or can emit errors for logging
type Event interface {
	Key() any
	Value() any
	Error() error

	// Event can be emitted to a channel. Returns false if unable to do so
	Emit(chan<- Event) bool
}

// RouterTask is a task which maps paths to routes
type RouterTask interface {
	// Add a prefix/path mapping to a handler for one or more HTTP methods
	AddHandler(prefix string, path *regexp.Regexp, fn http.HandlerFunc, methods ...string) error

	// Add middleware for a unique name
	AddMiddleware(name string, fn func(http.HandlerFunc) http.HandlerFunc) error

	// Set middleware for a prefix. Called from left to right.
	SetMiddleware(prefix string, chain ...string) error
}
