package nginx

import (
	"fmt"
	"io/fs"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type nginx struct {
	label     string
	available *Folder
	enabled   *Folder
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewWithConfig(fs fs.StatFS, c Config) (Task, error) {
	r := new(nginx)
	r.label = c.Label

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
	str += fmt.Sprintf(" available=%v", r.available)
	str += fmt.Sprintf(" enabled=%v", r.enabled)
	return str + ">"
}
