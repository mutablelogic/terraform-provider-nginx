package nginx

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	"github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type File struct {
	path    string
	info    fs.FileInfo
	data    []byte
	enabled string
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewFile(path string, info fs.FileInfo) *File {
	this := new(File)
	this.path = path
	this.info = info

	// Return success
	return this
}

func CreateFile(path string, data []byte) (*File, error) {
	if err := ioutil.WriteFile(path, data, defaultFileMode); err != nil {
		return nil, err
	}

	// Get file information
	info, err := os.Stat(path)
	if err != nil {
		os.Remove(path)
		return nil, err
	}

	// Create the file and set the data
	file := NewFile(path, info)
	file.data = data

	// Return success
	return file, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f *File) String() string {
	str := "<nginx-file"
	str += fmt.Sprintf(" name=%q", f.Name())
	str += fmt.Sprintf(" size=%d", f.info.Size())
	if f.enabled != "" {
		str += " enabled"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the path for the configuration
func (f *File) Path() string {
	return f.path
}

// Return the name of the configuration
func (f *File) Name() string {
	return strings.TrimSuffix(f.info.Name(), defaultExt)
}

// Set the enabled path
func (f *File) SetEnabled(path string) {
	f.enabled = path
}

// Return true if the configuration is enabled
func (f *File) Enabled() bool {
	return f.enabled != ""
}

// Read the configuration file
func (f *File) Read() ([]byte, error) {
	// If file has not changed, return the cached version
	if info, err := os.Stat(f.path); err != nil {
		return nil, err
	} else if info.ModTime() == f.info.ModTime() && f.data != nil {
		return f.data, nil
	} else {
		f.info = info
	}

	// Read the file
	if data, err := ioutil.ReadFile(f.path); err != nil {
		return nil, err
	} else {
		f.data = data
	}

	// Return success
	return f.data, nil
}

// Disable an enabled configuration
func (f *File) Disable() error {
	if f.enabled == "" {
		return ErrOutOfOrder
	}
	if err := os.Remove(f.enabled); err != nil {
		return err
	} else {
		f.enabled = ""
	}
	// Return success
	return nil
}

// Delete the file
func (f *File) Revoke() error {
	var result error
	if f.enabled != "" {
		if err := os.Remove(f.enabled); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if err := os.Remove(f.path); err != nil {
		result = multierror.Append(result, err)
	}

	// Blank fields
	f.info = nil
	f.data = nil
	f.enabled = ""

	// Return any errors
	return result
}
