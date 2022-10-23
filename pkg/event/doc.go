/*
Package `event` provides events and subscriptions. Event messages are created with
a key/value pair using `NewEvent`, or an error with `NewError`. Events can be
emitted on a channel using the `Emit` method, which returns false if the channel
is full.

In order to make an existing instance able to emit events, embed a `PubSub` instance
within the instance, and set the `Cap` value to the capacity required. For example,

	type Task struct {
		event.PubSub
	}

	func NewTaskWithCapacity(cap int) Task {
		task := new(Task)
		task.Cap = cap
		return task
	}
*/
package event
