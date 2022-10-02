package provider

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync"

	// Module imports
	multierror "github.com/hashicorp/go-multierror"
	iface "github.com/mutablelogic/terraform-provider-nginx"
	event "github.com/mutablelogic/terraform-provider-nginx/pkg/event"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type provider struct {
	event.PubSub

	// Enumeration of task plugins, keyed by name
	plugins map[string]reflect.Type

	// Enumeration of tasks, keyed by label
	tasks map[string]iface.Task
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	reTaskName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]+$`)
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// New creates a new empty provider with no tasks
func New() *provider {
	p := new(provider)
	p.plugins = make(map[string]reflect.Type)
	p.tasks = make(map[string]iface.Task)
	return p
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS - TASK

func (p *provider) Label() string {
	return "provider"
}

func (p *provider) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	var result error

	// Run all tasks
	for label, task := range p.tasks {
		wg.Add(2)

		// Emit events from task
		go func(task iface.Task) {
			defer wg.Done()
			ch := task.C()
			if ch != nil {
				for {
					select {
					case <-ctx.Done():
						return
					case event := <-ch:
						if event != nil && !p.Emit(event) {
							panic(fmt.Sprintln("Unable to emit: ", event))
						}
					}
				}
			}
		}(task)

		// Run task
		go func(label string, task iface.Task) {
			defer wg.Done()
			if err := task.Run(ctx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				result = multierror.Append(result, fmt.Errorf("%v: %w", label, err))
			}
		}(label, task)
	}

	// Wait until all tasks are completed
	wg.Wait()

	// TODO: Close tasks in the reverse order they were created

	// Close channel
	p.Emit(nil)

	// Return any errors
	return result
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *provider) String() string {
	str := "<provider"
	str += fmt.Sprintf(" label=%q", p.Label())
	return str + ">"
}

// New creates a new task from a configuration with a unique label
func (p *provider) New(ctx context.Context, config iface.TaskPlugin) (iface.Task, error) {
	name := config.Name()

	// Check the plugin by type
	t := reflect.TypeOf(config)
	if !reTaskName.MatchString(name) {
		return nil, ErrBadParameter.Withf("Invalid name %q for plugin", name)
	} else if t_, exists := p.plugins[name]; exists {
		if t != t_ {
			return nil, ErrDuplicateEntry.Withf("Plugin %q already exists", name)
		}
	} else {
		p.plugins[name] = t
	}

	// Create a new task
	task, err := config.New(ctx, p)
	if err != nil {
		return nil, err
	} else if task == nil {
		return nil, ErrInternalAppError.Withf("Unexpected nil return when creating task %q ", name)
	} else if label := task.Label(); !reTaskName.MatchString(label) {
		return nil, ErrBadParameter.Withf("Invalid label %q for task %q ", label, name)
	} else if _, exists := p.tasks[name+"."+label]; exists {
		return nil, ErrDuplicateEntry.Withf("Task with label %q already exists", name+"."+label)
	} else {
		p.tasks[name+"."+label] = task
	}

	// Return success
	return task, nil
}
