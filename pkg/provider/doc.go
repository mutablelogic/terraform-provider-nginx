//

/*
Package provider co-ordinates tasks and events for the lifecycle
of the application.

# Providers

Creating an empty provider is simple using the New function. Then, tasks are created using their
configurations and added to the provider.

	provider := provider.New()
	if task, err := provider.New(context.Background(), provider.Task{ Label: "task_instance_1" }); err != nil {
		// Handle error
	} else {
		// Use task
	}

Here, context is used to pass in additional values from the calling application. Once one more more
tasks have been created, they can be run:

	  if err := provider.Run(ctx); err != nil {
		// ...
	  }

This will run all the tasks in the background, ensuring that any dependencies for tasks are satisfied
when the context is cancelled.

# Events

Events are used to communicate between tasks. You can subscribe to the stream of events....
*/
package provider
