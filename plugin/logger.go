package plugin

import (
	"context"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

// Logger provides logging services
type Logger interface {
	Task

	// Log a message with arguments
	Log(context.Context, ...any)

	// Log a message with formatted arguments
	Logf(context.Context, string, ...any)
}
