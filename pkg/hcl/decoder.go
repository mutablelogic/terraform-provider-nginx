package hcl

import (
	"fmt"
	"io/fs"
	"reflect"

	// Module imports

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type decoder struct {
	tag     string
	spec    hcldec.TupleSpec
	blocks  []string
	plugins map[string]TaskPlugin
	vars    map[string]cty.Value
	fns     map[string]function.Function
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewDecoder() *decoder {
	d := new(decoder)

	// Add in a specification for var definitions
	d.tag = TagName
	d.plugins = make(map[string]TaskPlugin)
	d.vars = make(map[string]cty.Value)
	d.fns = make(map[string]function.Function)

	// Register the var block type
	d.MustRegister(varblock{})

	// Return success
	return d
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Register a block for decoding
func (d *decoder) MustRegister(plugin TaskPlugin) {
	if err := d.Register(plugin); err != nil {
		panic(err)
	}
}

// Register a block for decoding
func (d *decoder) Register(plugin TaskPlugin) error {
	// Register plugin
	name := plugin.Name()
	if _, exists := d.plugins[name]; exists {
		return ErrDuplicateEntry.Withf("%q", name)
	}

	// Append spec
	if spec, err := SpecForType(reflect.TypeOf(plugin), d.tag, name); err != nil {
		return err
	} else {
		d.blocks = append(d.blocks, name)
		d.spec = append(d.spec, spec)
		d.plugins[name] = plugin
	}

	// Return success
	return nil
}

// Parse a file
func (d *decoder) Parse(filesys fs.FS, path string) ([]TaskPlugin, error) {
	body, err := Parse(filesys, path)
	if err != nil {
		return nil, err
	}

	// Create references for variables in the body
	refs := NewRefs()
	if err := refs.CreateReferences(body, d.spec); err != nil {
		return nil, err
	}

	// Decode the body
	var plugins []TaskPlugin
	var result error
	if value, diags := hcldec.Decode(body, d.spec, refs.Context()); diags.HasErrors() {
		return nil, diags
	} else {
		// Convert cty.Value into configuration objects
		value.ForEachElement(func(key, tuple cty.Value) bool {
			// Obtain the index of the block
			var i int
			if err := FromCtyValue(key, &i); err != nil {
				result = multierror.Append(result, err)
				return true
			} else if i < 0 || i >= len(d.blocks) {
				result = multierror.Append(result, ErrInternalAppError)
				return true
			}

			// Obtain the plugin
			proto, exists := d.plugins[d.blocks[i]]
			if !exists {
				result = multierror.Append(result, ErrInternalAppError)
				return true
			}

			return tuple.ForEachElement(func(_, tuple cty.Value) bool {
				// Tuple should be iteratable for the fields
				plugin := fromPrototype(proto)
				if err := FromCtyValue(tuple, plugin); err != nil {
					result = multierror.Append(result, fmt.Errorf("%w: %q", err, plugin.Name()))
					return true
				} else {
					plugins = append(plugins, plugin)
				}
				return false
			})
		})
	}

	// Return success
	return plugins, result
}

/////////////////////////////////////////////////////////////////////
// CREATE OBJECT FROM PROTOTYPE

// fromPrototype creates a new zero-valued bject from a prototype
func fromPrototype(proto TaskPlugin) TaskPlugin {
	return reflect.New(reflect.TypeOf(proto)).Interface().(TaskPlugin)
}
