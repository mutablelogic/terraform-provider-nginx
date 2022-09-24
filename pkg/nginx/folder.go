package nginx

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Folder tracks files within a folder
type Folder struct {
	fs        fs.StatFS
	path      string
	recursive bool
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewFolder(fs fs.StatFS, path string, recursive bool) (*Folder, error) {
	f := new(Folder)

	// Remove initial '/' from path
	if strings.HasPrefix(path, pathSeparator) {
		path = path[1:]
	}

	// Check to make sure path is valid
	if info, err := fs.Stat(path); os.IsNotExist(err) {
		return nil, ErrNotFound.With(path)
	} else if err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, ErrBadParameter.With(path)
	} else {
		f.fs = fs
		f.path = path
		f.recursive = recursive
	}

	// Return success
	return f, nil
}

///////////////////////////////////////////////////////////////////////////////
// STRINFIGY

func (f *Folder) String() string {
	str := "<nginx-folder"
	str += fmt.Sprintf(" path=%q", f.path)
	if f.recursive {
		str += " recursive"
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (f *Folder) RelPath(root string) string {
	if root == "" {
		return f.path
	} else if rel, err := filepath.Rel(root, f.path); err != nil {
		return f.path
	} else {
		return rel
	}
}

func (f *Folder) Enumerate() ([]*File, error) {
	return enumerateFiles(f.fs, f.path, f.recursive)
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func enumerateFiles(filesys fs.StatFS, root string, recursive bool) ([]*File, error) {
	var result []*File

	if err := fs.WalkDir(filesys, root, func(path string, d fs.DirEntry, err error) error {
		// Skip errors
		if err != nil {
			return err
		}

		// Ignore hidden files
		if strings.HasPrefix(d.Name(), ".") && root != path {
			if d.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		// Recurse into directories
		if d.IsDir() && root != path {
			if recursive {
				return nil
			} else {
				return filepath.SkipDir
			}
		}

		// Only enumerate regular files
		if info, err := d.Info(); err != nil {
			return nil
		} else if validFileMode(info.Mode()) {
			result = append(result, NewFile(path, info))
		}

		// Return success
		return nil
	}); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func validFileMode(mode fs.FileMode) bool {
	if mode.IsRegular() {
		return true
	}
	if mode.Type() == fs.ModeSymlink {
		return true
	}
	return false
}
