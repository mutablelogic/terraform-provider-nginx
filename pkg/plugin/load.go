package plugin

import (
	"path/filepath"

	// Modules
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

// LoadPluginsForPattern will load and return a map of plugins for a given glob pattern,
// keyed against the plugin name.
func LoadPluginsForPattern(pattern string) (map[string]TaskPlugin, error) {
	var result error

	// Seek plugins
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Load plugins
	plugins := make(map[string]TaskPlugin, len(files))
	for _, path := range files {
		plugin, err := PluginWithPath(path)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		// Check for duplicate plugins
		name := plugin.Name()
		if _, exists := plugins[name]; exists {
			result = multierror.Append(result, ErrInternalAppError.Withf("Duplicate plugin: %q", name))
			continue
		}

		// Set plugin
		plugins[name] = plugin
	}

	// Return any errors
	return plugins, result
}
