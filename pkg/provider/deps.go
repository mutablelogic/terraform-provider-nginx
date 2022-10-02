package provider

import (
	"fmt"
	"reflect"

	// Modules
	iface "github.com/mutablelogic/terraform-provider-nginx"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	typeTaskPluginInterface = reflect.TypeOf((*iface.TaskPlugin)(nil)).Elem()
	typeTaskInterface       = reflect.TypeOf((*iface.Task)(nil)).Elem()
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// ReorderPlugins reorders the tasks so that they are in the correct order for
// cancellation
func ReorderPlugins(plugins ...iface.TaskPlugin) ([]iface.TaskPlugin, error) {
	result := []iface.TaskPlugin{}
	for _, plugin := range plugins {
		if deps := dependencies(plugin); deps == nil {
			return nil, ErrBadParameter.Withf("Invalid plugin: %q", plugin.Name())
		} else {
			result = append(result, deps...)
			fmt.Println(plugin, "=>", deps)
		}
	}

	// Return success
	return result, nil
}

// dependencies returns any field values that are referenced
func dependencies(plugin iface.TaskPlugin) []iface.TaskPlugin {
	// Check for parameters
	if plugin == nil {
		return nil
	}
	// dereference pointer
	rv := reflect.ValueOf(plugin)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	// Check for implementation of Task interface
	if rv.Kind() != reflect.Struct {
		return nil
	} else if !rv.Type().Implements(typeTaskPluginInterface) {
		return nil
	}
	return dependencies_(rv)
}

func dependencies_(rv reflect.Value) []iface.TaskPlugin {
	result := []iface.TaskPlugin{}
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		if !fv.Type().Implements(typeTaskInterface) {
			continue
		}
		// Field element is a task
		fmt.Println("Field", rv.Type().Field(i).Name)
	}
	return result
}
