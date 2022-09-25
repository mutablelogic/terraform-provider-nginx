package nginx

import (
	"context"
	"os"
	"path/filepath"

	// Modules
	util "github.com/mutablelogic/terraform-provider-nginx/pkg/util"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/terraform-provider-nginx"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Config struct {
	Label     string `hcl:"label,label"`                  // Label for the configuration
	Path      string `hcl:"conf_path"`                    // Root path for the configuration
	PidPath   string `hcl:"pid_path,optional"`            // Path to the PID file
	Available string `hcl:"available_path,optional"`      // Path to available sites, under root
	Recursive bool   `hcl:"available_recursive,optional"` // Recursively search in available folder
	Enabled   string `hcl:"enabled_path,optional"`        // Path to enabled sites, under root
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabel     = "nginx"
	defaultAvailable = "sites-available"
	defaultEnabled   = "sites-enabled"
	defaultPidPath   = "/run/nginx.pid"
	defaultExt       = ".conf"
	defaultFileMode  = 0644
	pathSeparator    = string(os.PathSeparator)
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (c Config) Name() string {
	return DefaultLabel
}

// Return a new task. Label for the task can be retrieved from context
func (c Config) New(ctx context.Context, provider Provider) (Task, error) {
	// Set label
	if c.Label == "" {
		c.Label = DefaultLabel
	}

	// Set PID path
	if c.PidPath == "" {
		c.PidPath = defaultPidPath
	}

	// Check label
	if !util.IsIdentifier(c.Label) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label)
	}

	// Check path
	if c.Path == "" || !filepath.IsAbs(c.Path) {
		if cwd, err := os.Getwd(); err != nil {
			return nil, ErrBadParameter.With(err)
		} else {
			c.Path = filepath.Join(cwd, c.Path)
		}
	}
	if info, err := os.Stat(c.Path); err != nil {
		return nil, ErrBadParameter.With(err)
	} else if !info.IsDir() {
		return nil, ErrBadParameter.With(c.Path)
	}

	// Set available path
	if c.Available == "" {
		c.Available = defaultAvailable
	}
	if !filepath.IsAbs(c.Available) {
		c.Available = filepath.Join(c.Path, c.Available)
	}

	// Set enabled path
	if c.Enabled == "" {
		c.Enabled = defaultEnabled
	}
	if !filepath.IsAbs(c.Enabled) {
		c.Enabled = filepath.Join(c.Path, c.Enabled)
	}

	// Return configuration
	return NewWithConfig(c)
}
