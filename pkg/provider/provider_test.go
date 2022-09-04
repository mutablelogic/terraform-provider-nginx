package provider_test

import (
	"context"
	"errors"
	"testing"

	// Module imports

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Provider_001(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("Unexpected nil returned from New")
	}
}

func Test_Provider_002(t *testing.T) {
	provider := New()
	if err := provider.Register(task{}); err != nil {
		t.Error(err)
	}
	if err := provider.Register(task{}); !errors.Is(err, ErrDuplicateEntry) {
		t.Error("Unexpected error:", err)
	}
}

func Test_Provider_003(t *testing.T) {
	provider := New()
	if err := provider.Register(task{}); err != nil {
		t.Error(err)
	}
	if task_, err := provider.New(context.Background(), task{}); err != nil {
		t.Error(err)
	} else if task_ == nil {
		t.Error("Unexpected nil task")
	} else if task__ := provider.TaskWithLabel(context.Background(), "label"); task__ != task_ {
		t.Error("Unexpected task returned from TaskWithLabel", task__)
	} else if task__ := provider.TasksWithName(context.Background(), "test"); !task_.(*task).in(task__) {
		t.Error("Unexpected task returned from TasksWithName", task_)
	}
}

func Test_Provider_004(t *testing.T) {
	provider := New()
	if err := provider.Register(task{}); err != nil {
		t.Error(err)
	}
	if _, err := provider.New(context.Background(), task{}); err != nil {
		t.Error(err)
	} else if _, err := provider.New(context.Background(), task{}); !errors.Is(err, ErrDuplicateEntry) {
		t.Error("Unexpected error:", err)
	}
}

/////////////////////////////////////////////////////////////////////
// TASK IMPLEMENTATION

type task struct {
	label string
}

func (t task) Name() string {
	return "test"
}

func (t task) New(ctx context.Context, provider Provider) (Task, error) {
	this := &t
	this.label = "label"
	return &t, nil
}

func (t *task) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (t *task) Label() string {
	return t.label
}

func (t *task) C() <-chan Event {
	return nil
}

func (t *task) in(tasks []Task) bool {
	for _, elem := range tasks {
		if t == elem {
			return true
		}
	}
	return false
}
