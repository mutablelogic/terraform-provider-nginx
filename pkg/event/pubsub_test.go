package event_test

import (
	"sync"
	"testing"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/event"
)

func Test_PubSub_001(t *testing.T) {
	p := new(PubSub)
	for i := 0; i < 100; i++ {
		if done := p.Emit(NewEvent(t.Name(), i)); !done {
			t.Fatal("Expected done")
		}
	}
	p.Emit(nil)
}

func Test_PubSub_002(t *testing.T) {
	var wg sync.WaitGroup

	p := new(PubSub)
	ch := p.Sub()
	n := 100
	m := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		for evt := range ch {
			t.Log("Received", evt)
			m += 1
		}
	}()

	for i := 0; i < n; i++ {
		if done := p.Emit(NewEvent(t.Name(), i)); !done {
			t.Fatal("Expected done")
		}
	}

	// Close channels
	p.Emit(nil)

	// Wait until goroutine completed
	wg.Wait()

	// Check number of received events
	if n != m {
		t.Error("Expected", n, "events, got", m)
	}
}

func Test_PubSub_003(t *testing.T) {
	var wg sync.WaitGroup

	p := new(PubSub)
	n := 100
	m := 0

	for r := 0; r < 5; r++ {
		ch := p.Sub()
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for evt := range ch {
				t.Log("Received", evt, "on channel", i)
				m += 1
			}
		}(r)
	}

	// Emit events in foreground
	for i := 0; i < n; i++ {
		if done := p.Emit(NewEvent(t.Name(), i)); !done {
			t.Fatal("Expected done")
		}
	}

	// Close channels
	p.Emit(nil)

	// Wait until goroutine completed
	wg.Wait()

	// Check number of received events
	if n*5 != m {
		t.Error("Expected", n, "events, got", m)
	}
}
