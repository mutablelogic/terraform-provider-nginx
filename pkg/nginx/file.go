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
	path string
	info fs.FileInfo
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
	str := "<file"
	str += fmt.Sprintf(" path=%q", f.path)
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (f *File) Path() string {
	return f.path
}
