package provider

import (
	"context"
	"fmt"
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
	// Enumeration of task plugins, keyed by label
	plugins map[string]TaskPlugin

	// Enumeration of tasks, keyed by label
	tasks map[string]task_
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
	p.plugins = make(map[string]TaskPlugin)
	p.tasks = make(map[string]task_)
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

	// Create tasks
	// TODO: These should be done in the right order
	//for label, task := range p.tasks {
	//	task, err := task.New()
	//}

	// TODO: Collect events

	// Run all tasks
	for label, task := range p.tasks {
		wg.Add(1)
		go func(label string, task Task) {
			defer wg.Done()
			if err := task.Run(ctx); err != nil {
				result = multierror.Append(result, fmt.Errorf("%v: %w", label, err))
			}
		}(label, task.Task)
	}

	// Wait until all tasks are completed
	wg.Wait()

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

/*
// Register a task with the provider
func (p *provider) Register(config TaskPlugin) error {
	name := config.Name()
	if !reTaskName.MatchString(name) {
		return ErrBadParameter.With(name)
	} else if _, exists := p.plugins[name]; exists {
		return ErrDuplicateEntry.With(name)
	} else {
		p.plugins[name] = config
	}

	// Return success
	return nil
}
*/

// New creates a new task from a configuration with a unique label
func (p *provider) New(ctx context.Context, config TaskPlugin) (Task, error) {
	// Create the task
	name := config.Name()
	task, err := config.New(ctx, p)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrInternalAppError.Withf("Unexpected nil return when creating task %q ", name)
	}
	if label := task.Label(); !reTaskName.MatchString(label) {
		return nil, ErrBadParameter.Withf("Invalid label %q for task %q ", label, name)
	} else if _, exists := p.tasks[label]; exists {
		return nil, ErrDuplicateEntry.Withf("Task %q with label %q already exists", name, label)
	} else {
		p.tasks[label] = task_{task, name, label}
	}

	// Return success
	return task, nil
}

// Return a task with the given label
// TODO: if called within New then creates a dependency between tasks on Run
func (p *provider) TaskWithLabel(ctx context.Context, label string) Task {
	return p.tasks[label].Task
}

// Return tasks with the given name
// TODO: if called within New then creates a dependency between tasks on Run
func (p *provider) TasksWithName(ctx context.Context, name string) []Task {
	result := []Task{}
	for _, task := range p.tasks {
		if name == task.name {
			result = append(result, task.Task)
		}
	}
	return result
}
