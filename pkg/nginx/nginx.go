package nginx

import (
	"fmt"
	"io/fs"

	// Namespace imports
	"github.com/hashicorp/go-multierror"
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

func NewWithConfig(fs fs.StatFS, c Config) (Task, error) {
	r := new(nginx)
	r.label = c.Label
	r.root = c.Path

	// Set up available folder
	if folder, err := NewFolder(fs, c.Available, c.Recursive); err != nil {
		return nil, err
	} else {
		r.available = folder
	}

	// Set up enabled folder
	if folder, err := NewFolder(fs, c.Enabled, false); err != nil {
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
		str += fmt.Sprintf(" path=%q", pathSeparator+r.root)
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

	fmt.Println(available, enabled)

	// Return success
	return nil, nil
}
