// config package manages the configuration of nginx
package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	// Module imports
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	AvailablePath string
	EnabledPath   string
}

type config struct {
	readonly     bool
	available    string
	enabled      string
	atime, etime time.Time
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	FiLE_CREATE_MODE = 0644
	DIR_CREATE_MODE  = 0755
)

var (
	reValidName = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) New() (*config, error) {
	this := new(config)
	// Create folders if they don't exist
	if path, err := checkPath(c.AvailablePath); err != nil {
		return nil, err
	} else {
		this.available = path
		if writable, err := util.IsWritableDir(path); err != nil {
			return nil, err
		} else {
			this.readonly = !writable
		}
	}
	if path, err := checkPath(c.EnabledPath); err != nil {
		return nil, err
	} else {
		this.enabled = path
		if writable, err := util.IsWritableDir(path); err != nil {
			return nil, err
		} else if !this.readonly {
			this.readonly = !writable
		}
	}

	// Return success
	return this, nil
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c *config) String() string {
	str := "<config"
	if c.available != "" {
		str += fmt.Sprintf(" available_path=%q", c.available)
	}
	if c.enabled != "" {
		str += fmt.Sprintf(" enabled_path=%q", c.enabled)
	}
	if c.readonly {
		str += " readonly"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// ENUMERATE CONFIGURATIONS

// EnumerateAvailable returns all files known to the configuration. Does not
// follow symlinks.
func (c *config) EnumerateAvailable() ([]File, error) {
	return enumerateFiles(c.available)
}

// EnumerateEnabled returns enabled files known to the configuration. Does
// follow symlinks for files
func (c *config) EnumerateEnabled() ([]File, error) {
	return enumerateFiles(c.enabled)
}

// Create will create a new configuration file, but not enable it. It can
// overwrite an existing file
func (c *config) Create(name string, body []byte) (File, error) {
	// Check for readonly
	if c.readonly {
		return nil, ErrOutOfOrder.With("readonly mode")
	}
	// Check for a valid configuration name
	name = strings.TrimSuffix(name, FILE_EXT)
	if !reValidName.MatchString(name) {
		return nil, ErrBadParameter.Withf("Create: %q", name)
	}
	// Create a new file
	path := filepath.Join(c.available, name+FILE_EXT)
	if err := os.WriteFile(path, body, FiLE_CREATE_MODE); err != nil {
		return nil, err
	} else if info, err := os.Stat(path); err != nil {
		os.Remove(path)
		return nil, err
	} else {
		return NewFile(path, info), nil
	}
}

// Remove will remove a configuration file
func (c *config) Remove(f File) error {
	// Preconditions
	if f == nil {
		return ErrBadParameter.With("Remove")
	}
	if c.readonly {
		return ErrOutOfOrder.With("Remove", "readonly mode")
	}
	// Unlink the file if it is enabled
	if err := c.Unlink(f); err != nil {
		return err
	}
	// Remove the file
	return os.Remove(f.Path())
}

// Enabled will return true if any file is linked into the enabled folder
func (c *config) Enabled(f File) (bool, error) {
	// Preconditions
	if f == nil {
		return false, ErrBadParameter.With("Enabled")
	}

	// Check if file is already enabled
	enabled_path := filepath.Join(c.enabled, f.Filename())
	if info, err := os.Stat(enabled_path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		enabled := NewFile(enabled_path, info)
		if matches, err := enabled.Matches(f); err != nil {
			return false, err
		} else {
			return matches, nil
		}
	}
}

// Link will create a new symlink from the enabled path to the available path
// Effectively enabling the configuration. It will return the linked file.
func (c *config) Link(f File) (File, error) {
	// Preconditions
	if f == nil {
		return nil, ErrBadParameter.With("Link")
	}
	if c.readonly {
		return nil, ErrOutOfOrder.With("readonly mode")
	}

	// Check if file is already enabled, ignore
	enabled_path := filepath.Join(c.enabled, f.Filename())
	if info, err := os.Stat(enabled_path); err == nil {
		enabled := NewFile(enabled_path, info)
		if matches, err := enabled.Matches(f); err != nil {
			return nil, err
		} else if matches {
			return enabled, nil
		} else if err := os.RemoveAll(enabled_path); err != nil {
			return nil, err
		}
	}

	// Perform the link
	if err := os.Symlink(f.Path(), enabled_path); err != nil {
		return nil, err
	} else if info, err := os.Stat(enabled_path); err != nil {
		os.RemoveAll(enabled_path)
		return nil, err
	} else {
		return NewFile(enabled_path, info), nil
	}
}

// Unlink will remove a symlink from the enabled path given file at available path
// Effectively disabling the configuration
func (c *config) Unlink(f File) error {
	// Preconditions
	if f == nil {
		return ErrBadParameter.With("Unlink")
	}
	if c.readonly {
		return ErrOutOfOrder.With("readonly mode")
	}

	// Check if file is enabled, and remove it if it matches
	enabled_path := filepath.Join(c.enabled, f.Filename())
	if info, err := os.Stat(enabled_path); err == nil {
		enabled := NewFile(enabled_path, info)
		if matches, err := enabled.Matches(f); err != nil {
			return err
		} else if !matches {
			return ErrOutOfOrder.With("enabled file does not match")
		} else if err := os.RemoveAll(enabled_path); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

// Changed returns true if either available or enabled folder have changed
func (c *config) Changed() (bool, error) {
	var changed bool

	// Check available
	if info, err := os.Stat(c.available); err != nil {
		return false, err
	} else if !c.atime.Equal(info.ModTime()) {
		changed = true
		c.atime = info.ModTime()
	}

	// Check enabled
	if info, err := os.Stat(c.enabled); err != nil {
		return false, err
	} else if !c.etime.Equal(info.ModTime()) {
		changed = true
		c.etime = info.ModTime()
	}

	// Return changed flag
	return changed, nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// check path exists, if not then try and create it
func checkPath(path string) (string, error) {
	if path_, err := filepath.Abs(path); err != nil {
		return "", err
	} else {
		path = path_
	}
	if info, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, DIR_CREATE_MODE); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else if !info.IsDir() {
		return "", ErrBadParameter.With(path)
	}

	// Return success
	return path, nil
}

func enumerateFiles(path string) ([]File, error) {
	var result []File

	if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		// Skip errors
		if err != nil {
			return err
		}

		// Ignore hidden files
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		// Recurse into directories
		if d.IsDir() {
			return nil
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
