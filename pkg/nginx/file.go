package nginx

import (
	"fmt"
	"io/fs"
	// Module imports
	// Namespace imports
	//. "github.com/djthorpe/go-errors"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type File struct {
	path      string
	info      fs.FileInfo
	data      []byte
	available bool
	enabled   bool
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

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f *File) String() string {
	str := "<nginx-file"
	str += fmt.Sprintf(" name=%q", f.Name())
	str += fmt.Sprintf(" size=%d", f.info.Size())
	if f.available {
		str += " available"
	}
	if f.enabled {
		str += " enabled"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (f *File) Path() string {
	return f.path
}

func (f *File) Name() string {
	return f.info.Name()
}

func (f *File) SetEnabled(v bool) {
	f.enabled = v
}

func (f *File) SetAvailable(v bool) {
	f.available = v
}

func (f *File) Enabled() bool {
	return f.enabled
}

func (f *File) Available() bool {
	return f.available
}

func (f *File) Read(filesys fs.FS) ([]byte, error) {
	// If file has not changed, return the cached version
	if info, err := filesys.(fs.StatFS).Stat(f.path); err != nil {
		return nil, err
	} else if info.ModTime() == f.info.ModTime() && f.data != nil {
		return f.data, nil
	} else {
		f.info = info
	}

	// Read the file
	if data, err := fs.ReadFile(filesys, f.path); err != nil {
		return nil, err
	} else {
		f.data = data
	}

	// Return success
	return f.data, nil
}
