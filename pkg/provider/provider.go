package provider

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync"

	// Module imports
	"github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type provider struct {
	// Enumeration of task plugins, keyed by name
	plugins map[string]reflect.Type

	// Enumeration of tasks, keyed by label
	tasks map[string]task_

	// Event channel
	ch chan Event
}

type task_ struct {
	Task

	// Name of task (not necessarily unique)
	name string

	// Label of the task (unique)
	label string
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	reTaskName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]+$`)
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func New() *provider {
	p := new(provider)
	p.plugins = make(map[string]reflect.Type)
	p.tasks = make(map[string]task_)
	p.ch = make(chan Event)
	return p
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS - TASK

// Return channels for events
func (p *provider) C() <-chan Event {
	return nil
}

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
		go func(task Task) {
			defer wg.Done()
			ch := task.C()
			if ch != nil {
				for {
					select {
					case <-ctx.Done():
						return
					case event := <-ch:
						if !event.Emit(p.ch) {
							panic(fmt.Sprint("Unable to emit:", event))
						}
					}
				}
			}
		}(task)

		// Run task
		go func(label string, task Task) {
			defer wg.Done()
			if err := task.Run(ctx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				result = multierror.Append(result, fmt.Errorf("%v: %w", label, err))
			}
		}(label, task.Task)
	}

	// Wait until all tasks are completed
	wg.Wait()

	// TODO: Close tasks in the reverse order they were created

	// Close channel
	close(p.ch)

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
func (p *provider) New(ctx context.Context, config TaskPlugin) (Task, error) {
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
	} else if _, exists := p.tasks[label]; exists {
		return nil, ErrDuplicateEntry.Withf("Task %q with label %q already exists", name, label)
	} else {
		p.tasks[label] = task_{task, name, label}
	}

	// Return success
	return task, nil
}

// TaskWithLabel return a task with the given label or nil if not found
func (p *provider) TaskWithLabel(label string) Task {
	return p.tasks[label].Task
}

// TasksWithName returns a slice of tasks with the given name, or if name is empty
// return all tasks
func (p *provider) TasksWithName(name string) []Task {
	result := []Task{}
	for _, task := range p.tasks {
		if name == task.name || name == "" {
			result = append(result, task.Task)
		}
	}
	return result
}
