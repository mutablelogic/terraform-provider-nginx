package kernel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	// Import modules
	multierror "github.com/hashicorp/go-multierror"

	// Import interfaces
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	"github.com/mutablelogic/terraform-provider-nginx/pkg/event"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type kernel struct {
	RouterTask
	wg    sync.WaitGroup
	tasks map[string]Task
	ch    chan Event
}

type KernelEvent uint

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	KernelEventNone KernelEvent = iota
	KernelEventStart
	KernelEventStop
	KernelEventError
)

var (
	reTaskName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]+$`)
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func New() *kernel {
	k := new(kernel)
	k.tasks = make(map[string]Task)
	k.ch = make(chan Event, 1000)
	return k
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v KernelEvent) String() string {
	switch v {
	case KernelEventStart:
		return "start"
	case KernelEventStop:
		return "stop"
	case KernelEventError:
		return "error"
	default:
		return "???"
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (k *kernel) Run(ctx context.Context) error {
	var result error
	var counter sync.WaitGroup

	// If there is no routertask, then quit
	if k.RouterTask == nil {
		return ErrNotFound.With("RouterTask")
	}

	// Run all the tasks in the background, except for the router
	counter.Add(len(k.tasks) - 1)
	for key, task := range k.tasks {
		// Do not start the router task until all the other tasks are started
		if task == k.RouterTask.(Task) {
			continue
		}

		// Start other tasks in the background
		k.wg.Add(1)
		counter.Done()
		go func(key string, task Task) {
			defer k.wg.Done()
			if err := k.startTask(ctx, key, task); err != nil {
				result = multierror.Append(result, err)
			}
		}(key, task)
	}

	// Wait for all tasks (except router task) to have been started
	time.Sleep(time.Second) // HACK
	counter.Wait()

	// If there are any errors, then return
	if result != nil {
		return result
	}

	// Start the router task
	for key, task := range k.tasks {
		// Do not start the router task until all the other tasks are started
		if task != k.RouterTask.(Task) {
			continue
		}

		// Start other tasks in the background
		k.wg.Add(1)
		go func(key string, task Task) {
			defer k.wg.Done()
			if err := k.startTask(ctx, key, task); err != nil {
				result = multierror.Append(result, err)
			}
		}(key, task)
	}

	// Wait for cancellation and all go-routines to complete
	<-ctx.Done()
	k.wg.Wait()

	// Close events channel
	close(k.ch)

	// Return any errors
	return result
}

func (k *kernel) Add(key string, task Task) error {
	// Precondition checks
	if !reTaskName.MatchString(key) {
		return ErrBadParameter.Withf("%q", key)
	}
	if task == nil {
		return ErrBadParameter.Withf("%q", key)
	}
	if _, exists := k.tasks[key]; exists {
		return ErrDuplicateEntry.Withf("%q", key)
	}

	// Set RouterTask
	if _, ok := task.(RouterTask); ok {
		if k.RouterTask != nil {
			return ErrDuplicateEntry.Withf("%q", key)
		} else {
			k.RouterTask = task.(RouterTask)
		}
	}

	// Set task
	k.tasks[key] = task

	// Create a goroutine to receive events from the task
	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		ch := task.C()
		if ch != nil {
			for evt := range task.C() {
				if !evt.Emit(k.ch) {
					fmt.Fprintln(os.Stderr, "kernel: event channel is full")
				}
			}
		}
	}()

	// Return success
	return nil
}

func (k *kernel) Get(key string) Task {
	return k.tasks[key]
}

func (k *kernel) C() <-chan Event {
	return k.ch
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (k *kernel) startTask(ctx context.Context, key string, task Task) error {
	// Start task and wait until done
	event.NewEvent(KernelEventStart, task).Emit(k.ch)
	if err := task.Run(ctx, k); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		event.NewError(fmt.Errorf("%v: %w", key, err)).Emit(k.ch)
		return err
	}

	// Return success
	event.NewEvent(KernelEventStop, key).Emit(k.ch)
	return nil
}
