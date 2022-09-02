package hcl

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	// Modules

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type config struct {
	path      string
	resources map[string]any
	variables map[string]cty.Value
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	fileExtHCL    = ".hcl"
	fileExtJSON   = ".json"
	pathSeparator = string(os.PathSeparator)
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// New creates an empty configuration from a working directory and a set of
// configuration "prototypes"
func New(root, path string, resources map[string]any) (*config, error) {
	this := new(config)
	this.resources = make(map[string]any)
	this.variables = make(map[string]cty.Value)

	// Set the path to the configuration, which removes any initial '/'
	if path := abspath(root, path); path == "" {
		return nil, ErrBadParameter.With("path")
	} else {
		this.path = path
	}

	// Create resource prototypes
	for name, resource := range resources {
		if _, exists := this.resources[name]; exists {
			return nil, ErrDuplicateEntry.With(name)
		} else if proto := prototypeOf(reflect.ValueOf(resource)); proto == nil {
			return nil, ErrInternalAppError.With(name)
		} else {
			this.resources[name] = proto
		}
	}

	// Return success
	return this, nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *config) String() string {
	str := "<hcl-config"
	str += fmt.Sprintf(" path=%q", pathSeparator+c.path)
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (c *config) Parse() error {
	parser := hclparse.NewParser()

	// Load HCL or JSON files
	filesys := os.DirFS("/")
	if err := fs.WalkDir(filesys, c.path, func(path string, info fs.DirEntry, err error) error {
		return c.walkconfig(parser, filesys, path, info, err)
	}); err != nil {
		return err
	}

	// Merge files together
	files := parser.Files()
	body := make([]*hcl.File, 0, len(files))
	for _, file := range files {
		body = append(body, file)
	}

	// Create a specification for parsing
	//	tuples := make([]string, 0, len(c.resources))
	spec := make(hcldec.TupleSpec, 0, len(c.resources))
	/*
		for name, resource := range c.resources {
			tuples = append(tuples, name)
			spec = append(spec, specOf(resource))
		}
	*/
	// Parse the configuration into a cty.Value
	value, diags := hcldec.Decode(hcl.MergeFiles(body), spec, &hcl.EvalContext{
		Variables: c.variables,
	})
	if diags.HasErrors() {
		return diags
	}

	// Convert cty.Value into configuration objects for plugins
	var result error
	value.ForEachElement(func(key, tuple cty.Value) bool {
		fmt.Println(key, "=>", tuple)
		return false
		/*
			var i int
			if err := gocty.FromCtyValue(key, &i); err != nil {
				result = multierror.Append(result, ErrInternalAppError)
				return true
			}
		*/
		/*

			// Get the plugin
			plugin, exists := c.Plugins[tuples[i]]
			if !exists {
				result = multierror.Append(result, ErrInternalAppError)
				return true
			}

			// Append plugin resources
			return tuple.ForEachElement(func(_, resource cty.Value) bool {
				if err := plugin.Append(resource); err != nil {
					result = multierror.Append(result, err)
					return true
				} else {
					return false
				}
			})

		*/
	})

	// Return success
	return result
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (c *config) walkconfig(parser *hclparse.Parser, filesys fs.FS, path string, d fs.DirEntry, err error) error {
	// Pass down any error
	if err != nil {
		return err
	}
	// Ignore any hidden files
	if strings.HasPrefix(d.Name(), ".") {
		if d.IsDir() {
			return fs.SkipDir
		} else {
			return nil
		}
	}
	// Recurse into folders
	if d.IsDir() {
		return nil
	}
	// Return error if not a regular file
	if !d.Type().IsRegular() {
		return ErrNotImplemented.Withf("%q", d.Name())
	}
	// Deal with filetypes
	switch strings.ToLower(filepath.Ext(d.Name())) {
	case fileExtHCL:
		if data, err := readAll(filesys, path); err != nil {
			return err
		} else if _, diags := parser.ParseHCL(data, path); diags.HasErrors() {
			return diags
		}
	case fileExtJSON:
		if data, err := readAll(filesys, path); err != nil {
			return err
		} else if _, diags := parser.ParseJSON(data, path); diags.HasErrors() {
			return diags
		}
	default:
		return ErrNotImplemented.Withf("%q", d.Name())
	}
	return nil
}

// Return an absolute path, but without the first '/'
func abspath(root, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	return strings.TrimPrefix(path, string(os.PathSeparator))
}

// Read all bytes from regular file
func readAll(filesys fs.FS, path string) ([]byte, error) {
	if data, err := fs.ReadFile(filesys, path); err != nil {
		return nil, fmt.Errorf("%q: %w", filepath.Base(path), err)
	} else {
		return data, nil
	}
}

// Return a reflect.Struct or nil if the prototype is not a struct
func prototypeOf(prototype reflect.Value) any {
	if prototype.Kind() == reflect.Ptr {
		prototype = prototype.Elem()
	}
	if prototype.Kind() != reflect.Struct {
		return nil
	} else if prototype.NumField() == 0 {
		return nil
	} else {
		return prototype.Interface()
	}
}
