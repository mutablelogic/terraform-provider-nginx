package provider_test

import (
	"context"
	"errors"
	"testing"
	"time"

	// Module imports

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"
)

/////////////////////////////////////////////////////////////////////
// TESTS

func Test_Provider_001(t *testing.T) {
	provider := New()
	if provider == nil {
		t.Fatal("Unexpected nil returned from New")
	}

	// Create a task
	if task, err := provider.New(context.Background(), Config{Label: "label"}); err != nil {
		t.Error(err)
	} else if task == nil {
		t.Error("Unexpected nil returned from New")
	} else if task.Label() != "label" {
		t.Error("Unexpected task label")
	}

	// Create a second task should return an error
	if _, err := provider.New(context.Background(), Config{Label: "label"}); !errors.Is(err, ErrDuplicateEntry) {
		t.Error("Unexpected error from New:", err)
	}

	/*
		// Task should be returned based on label
		if task := provider.TaskWithLabel("label"); task == nil {
			t.Error("Unexpected nil returned from TaskWithLabel")
		} else if task.Label() != "label" {
			t.Error("Unexpected task label:", task)
		}

		// One task returned based on name
		if tasks := provider.TasksWithName("test"); len(tasks) != 1 {
			t.Error("Unexpected nil returned from TasksWithName")
		}*/
}

func Test_Provider_002(t *testing.T) {
	provider := New()
	_, err := provider.New(context.Background(), Config{Label: "label"})
	if err != nil {
		t.Fatal(err)
	}
	// Creating a second task with the same name but different type should also fail
	_, err = provider.New(context.Background(), Config{Label: "label2"})
	if !errors.Is(err, ErrDuplicateEntry) {
		t.Fatal("Expected failure, got:", err)
	}
}

func Test_Provider_003(t *testing.T) {
	// Check label identifiers
	provider := New()
	_, err := provider.New(context.Background(), Config{Label: "00label"})
	if !errors.Is(err, ErrBadParameter) {
		t.Fatal("Expected failure, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label: "label 00"})
	if !errors.Is(err, ErrBadParameter) {
		t.Fatal("Expected failure, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label: "label.00"})
	if !errors.Is(err, ErrBadParameter) {
		t.Fatal("Expected failure, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label: "label-00"})
	if err != nil {
		t.Fatal("Expected success, got:", err)
	}
}

func Test_Provider_004(t *testing.T) {
	// Two task instances with different labels should be ok
	provider := New()
	_, err := provider.New(context.Background(), Config{Label: "task0"})
	if err != nil {
		t.Fatal("Expected success, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label: "task1"})
	if err != nil {
		t.Fatal("Expected success, got:", err)
	}
	/*
		if task := provider.TaskWithLabel("task0"); task == nil {
			t.Fatal("Expected task0, got nil")
		} else if task.Label() != "task0" {
			t.Fatal("Expected task0, got:", task.Label())
		}
		if task := provider.TaskWithLabel("task1"); task == nil {
			t.Fatal("Expected task1, got nil")
		} else if task.Label() != "task1" {
			t.Fatal("Expected task0, got:", task.Label())
		}*/
}

func Test_Provider_005(t *testing.T) {
	provider := New()
	// Create a task
	_, err := provider.New(context.Background(), Config{Label: "task0"})
	if err != nil {
		t.Fatal("Expected success, got:", err)
	}
	// Run task for 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := provider.Run(ctx); err != nil {
		t.Fatal("Expected success, got:", err)
	}
}
