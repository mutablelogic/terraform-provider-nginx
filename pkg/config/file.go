package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	// Module imports
	// Namespace imports
	//. "github.com/djthorpe/go-errors"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type file struct {
	path string
	name string
	info fs.FileInfo
	hash []byte
}

/////////////////////////////////////////////////////////////////////
// INTERFACES

type File interface {
	Path() string
	Name() string
	Filename() string
	Hash() (string, error)
	Matches(File) (bool, error)
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	FILE_EXT = ".conf"
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewFile(path string, info fs.FileInfo) *file {
	this := new(file)
	this.path = path
	this.info = info
	this.name = strings.TrimSuffix(this.info.Name(), FILE_EXT)

	// Return success
	return this
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f *file) String() string {
	str := "<file"
	str += fmt.Sprintf(" name=%q", f.Name())
	str += fmt.Sprintf(" filename=%q", f.Filename())
	if h, err := f.Hash(); err != nil {
		str += fmt.Sprint(" hash_error=", err)
	} else {
		str += fmt.Sprintf(" hash=%q", h)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (f *file) Path() string {
	return f.path
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Filename() string {
	return f.info.Name()
}

func (f *file) Hash() (string, error) {
	// If file has changed, then clear hash
	if info, err := os.Stat(f.path); err != nil {
		return "", err
	} else if !info.ModTime().Equal(f.info.ModTime()) {
		f.hash = nil
		f.info = info
	}

	if f.hash == nil {
		// Compute hash
		fh, err := os.Open(f.path)
		if err != nil {
			return "", err
		}
		defer fh.Close()
		hash := md5.New()
		if _, err := io.Copy(hash, fh); err != nil {
			return "", err
		} else {
			f.hash = hash.Sum(nil)
		}
	}
	return hex.EncodeToString(f.hash), nil
}

func (f *file) Matches(o File) (bool, error) {
	if h1, err := f.Hash(); err != nil {
		return false, err
	} else if h2, err := o.Hash(); err != nil {
		return false, err
	} else {
		return h1 == h2, nil
	}
}
