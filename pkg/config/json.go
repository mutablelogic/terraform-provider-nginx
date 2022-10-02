package config

import (
	"encoding/json"
	"io/fs"
	"reflect"

	// Modules
	multierror "github.com/hashicorp/go-multierror"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Resource struct {
	Name string `json:"resource"`
	Path string `json:"-"`
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// LoadJSONForPattern returns resources for a pattern of JSON files
func LoadJSONForPattern(filesys fs.FS, pattern string) ([]Resource, error) {
	var result error
	resources := []Resource{}

	// Find files
	files, err := fs.Glob(filesys, pattern)
	if err != nil {
		return nil, err
	}

	// Parse each file for resource definition
	for _, path := range files {
		r := Resource{Path: path}
		if data, err := fs.ReadFile(filesys, path); err != nil {
			result = multierror.Append(result, err)
		} else if err := json.Unmarshal(data, &r); err != nil {
			result = multierror.Append(result, err)
		} else if !util.IsIdentifier(r.Name) {
			result = multierror.Append(result, ErrBadParameter.Withf("Invalid resource: %q", path))
		} else {
			resources = append(resources, r)
		}
	}

	// Return result
	return resources, result
}

func ParseJSONResource(filesys fs.FS, resource Resource, plugin TaskPlugin) (TaskPlugin, error) {
	plugin = newPluginInstance(plugin)
	if data, err := fs.ReadFile(filesys, resource.Path); err != nil {
		return nil, err
	} else if err := json.Unmarshal(data, plugin); err != nil {
		return nil, err
	}

	// Return success
	return plugin, nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func newPluginInstance(plugin TaskPlugin) TaskPlugin {
	rv := reflect.New(reflect.TypeOf(plugin))
	return rv.Interface().(TaskPlugin)
}
