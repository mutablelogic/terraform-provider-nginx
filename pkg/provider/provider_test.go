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
	if task, err := provider.New(context.Background(), Config{Label_: "label"}); err != nil {
		t.Error(err)
	} else if task == nil {
		t.Error("Unexpected nil returned from New")
	}

	// Create a second task should return an error
	if _, err := provider.New(context.Background(), Config{Label_: "label"}); !errors.Is(err, ErrDuplicateEntry) {
		t.Error("Unexpected error from New:", err)
	}
}

func Test_Provider_003(t *testing.T) {
	// Check label identifiers
	provider := New()
	_, err := provider.New(context.Background(), Config{Label_: "00label"})
	if !errors.Is(err, ErrBadParameter) {
		t.Fatal("Expected failure, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label_: "label 00"})
	if !errors.Is(err, ErrBadParameter) {
		t.Fatal("Expected failure, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label_: "label.00"})
	if !errors.Is(err, ErrBadParameter) {
		t.Fatal("Expected failure, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label_: "label-00"})
	if err != nil {
		t.Fatal("Expected success, got:", err)
	}
}

func Test_Provider_004(t *testing.T) {
	// Two task instances with different labels should be ok
	provider := New()
	_, err := provider.New(context.Background(), Config{Label_: "task0"})
	if err != nil {
		t.Fatal("Expected success, got:", err)
	}
	_, err = provider.New(context.Background(), Config{Label_: "task1"})
	if err != nil {
		t.Fatal("Expected success, got:", err)
	}
}

func Test_Provider_005(t *testing.T) {
	provider := New()
	// Create a task
	_, err := provider.New(context.Background(), Config{Label_: "task0"})
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
