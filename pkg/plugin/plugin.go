package plugin

import (
	"plugin"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	funcConfig = "Config"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// PluginWithPath returns a plugin from a file path
func PluginWithPath(path string) (*TaskPlugin, error) {
	// Create a new module from plugin
	if plugin, err := plugin.Open(path); err != nil {
		return nil, err
	} else if fn, err := plugin.Lookup(funcConfig); err != nil {
		return nil, err
	} else if fn_, ok := fn.(func() TaskPlugin); !ok {
		return nil, ErrInternalAppError.With("New returned nil: ", path)
	} else if config := fn_(); config == nil {
		return nil, ErrInternalAppError.With("New returned nil: ", path)
	} else {
		return &config, nil
	}
}
