package logger

import (
	"context"
	"fmt"
	"log"

	// Modules
	"github.com/mutablelogic/terraform-provider-nginx/pkg/provider"

	// Namespace imports
	//. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type logger struct {
	provider.Task
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (Logger, error) {
	r := new(logger)

	r.SetLabel(c.Label)

	// Return success
	return r, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r *logger) String() string {
	str := "<logger"
	str += fmt.Sprintf(" label=%q", r.Label())
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Log a message with arguments
func (r *logger) Log(_ context.Context, args ...any) {
	log.Print(args...)
}

// Log a message with formatted arguments
func (r *logger) Logf(_ context.Context, format string, args ...any) {
	log.Printf(format, args...)
}
