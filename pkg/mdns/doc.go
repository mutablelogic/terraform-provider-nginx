/*
The mdns package listens and broadcasts for DNS packets

It listens on IP4 and IP6 connections for messages, and emits them
to any subscribers. In order to create a task instance,
use the `New` function with a `Config` object. The `Run` function
is then used to start the task instance, until the passed context is
cancelled, which closes the channels for all subscribers.

A `Message` can be created from a DNS packet (an answer) or with a
question, and sent using the `Send` method.
*/
package mdns
