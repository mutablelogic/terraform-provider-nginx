package provider

import (
	"context"
	"regexp"

	// Module imports
	environment "github.com/mutablelogic/terraform-provider-nginx/pkg/context"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type provider struct {
	// Enumeration of task plugins
	plugins map[string]TaskPlugin

	// Enumeration of tasks
	tasks map[string]task_
}

type task_ struct {
	Task

	// Name of the task
	name string
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

func (p *provider) Run(ctx context.Context) error {
	// Call new function on each task
	return nil
}

func (p *provider) Label() string {
	return ""
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p *provider) String() string {
	str := "<provider"
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

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

// New creates a new task from a configuration with a unique label
func (p *provider) New(ctx context.Context, config TaskPlugin) (Task, error) {
	// Create the task
	name := config.Name()
	task, err := config.New(environment.WithName(ctx, name), p)
	if err != nil {
		return nil, err
	} else if task == nil {
		return nil, ErrInternalAppError.Withf("Unexpected nil return when creating task %q ", name)
	} else if label := task.Label(); !reTaskName.MatchString(label) {
		return nil, ErrBadParameter.Withf("Invalid label %q for task %q ", label, name)
	} else if _, exists := p.tasks[label]; exists {
		return nil, ErrDuplicateEntry.Withf("Task %q with label %q already exists", name, label)
	} else {
		p.tasks[label] = task_{task, name}
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

// Return channels for events
func (p *provider) C() <-chan Event {
	return nil
}
