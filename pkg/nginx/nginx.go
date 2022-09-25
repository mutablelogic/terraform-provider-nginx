package nginx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// Modules
	multierror "github.com/hashicorp/go-multierror"
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
	. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type nginx struct {
	label     string
	root      string
	available *Folder
	enabled   *Folder
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(c Config) (Task, error) {
	r := new(nginx)
	r.label = c.Label
	r.root = c.Path

	// Set up available folder
	if folder, err := NewFolder(c.Available, c.Recursive); err != nil {
		return nil, err
	} else {
		r.available = folder
	}

	// Set up enabled folder
	if folder, err := NewFolder(c.Enabled, false); err != nil {
		return nil, err
	} else {
		r.enabled = folder
	}

	// Return success
	return r, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (r *nginx) String() string {
	str := "<nginx"
	str += fmt.Sprintf(" label=%q", r.label)
	if r.root != "" {
		str += fmt.Sprintf(" path=%q", r.root)
	}
	str += fmt.Sprintf(" available=%q", r.available.RelPath(r.root))
	str += fmt.Sprintf(" enabled=%q", r.enabled.RelPath(r.root))
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (r *nginx) Enumerate() ([]NginxConfig, error) {
	var result error

	// Enumerate available and enabled files
	available, err := r.available.Enumerate()
	if err != nil {
		result = multierror.Append(result, err)
	}

	enabled, err := r.enabled.Enumerate()
	if err != nil {
		result = multierror.Append(result, err)
	}
	// Return any errors
	if result != nil {
		return nil, result
	}

	// Create map of available files, based on MD5 hash
	config := make(map[string]*File, len(available))
	for _, file := range available {
		data, err := file.Read()
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}
		hash := util.MD5Hash(data)
		if _, exists := config[hash]; exists {
			result = multierror.Append(result, ErrInternalAppError.Withf("%v: duplicate", file.Path()))
			continue
		} else {
			config[hash] = file
		}
	}

	// If there are errors, return them
	if result != nil {
		return nil, result
	}

	// Set enabled if the file exists in the map
	for _, file := range enabled {
		data, err := file.Read()
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}
		hash := util.MD5Hash(data)
		if configfile, exists := config[hash]; !exists {
			config[hash] = file
			file.SetEnabled(configfile.Path())
		} else {
			configfile.SetEnabled(configfile.Path())
		}
	}

	// Create a set of configs
	configs := make([]NginxConfig, 0, len(config))
	for _, file := range config {
		configs = append(configs, file)
	}

	// Return success
	return configs, nil
}

// Enable a configuration
func (r *nginx) Enable(file NginxConfig) error {
	file_, ok := file.(*File)
	if !ok || file_ == nil {
		return ErrBadParameter
	}

	// Create a path for the file, based on the existing filename
	enabled_path := filepath.Join(r.enabled.RelPath(""), file_.Name())
	if err := os.Symlink(file_.Path(), enabled_path); err != nil {
		return err
	} else {
		file_.SetEnabled(enabled_path)
	}

	// Return success
	return nil
}

// Disable a configuration
func (r *nginx) Disable(file NginxConfig) error {
	file_, ok := file.(*File)
	if !ok || file_ == nil {
		return ErrBadParameter
	}
	return file_.Disable()
}

// Create a configuration
func (r *nginx) Create(name string, data []byte) (NginxConfig, error) {
	// Check parameters
	name = strings.TrimPrefix(name, defaultExt)
	if !util.IsIdentifier(name) {
		return nil, ErrBadParameter.Withf("Invalid name: %q", name)
	}
	if len(data) == 0 {
		return nil, ErrBadParameter.Withf("Invalid data")
	}

	// If path already exists, then error
	path := filepath.Join(r.available.RelPath(""), name+defaultExt)
	if _, err := os.Stat(path); err == nil {
		return nil, ErrDuplicateEntry.With(name)
	}

	// Create file and return it
	return CreateFile(path, data)
}

// Revoke a configuration
func (r *nginx) Revoke(file NginxConfig) error {
	file_, ok := file.(*File)
	if !ok || file_ == nil {
		return ErrBadParameter
	}
	return file_.Revoke()
}
