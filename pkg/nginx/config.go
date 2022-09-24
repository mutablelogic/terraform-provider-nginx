package nginx

import (
	"context"
	"io/fs"
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
	Label     string `hcl:"label,label"`        // Label for the configuration
	Path      string `hcl:"root"`               // Root path for the configuration
	Available string `hcl:"available,optional"` // Path to available sites, under root
	Recursive bool   `hcl:"recursive,optional"` // Recursively search in available folder
	Enabled   string `hcl:"enabled,optional"`   // Path to enabled sites, under root
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DefaultLabel     = "nginx"
	defaultAvailable = "sites-available"
	defaultEnabled   = "sites-enabled"
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

	// Check label
	if !util.IsIdentifier(c.Label) {
		return nil, ErrBadParameter.Withf("label: %q", c.Label)
	}

	// Set available and enabled if not set
	if c.Available == "" {
		c.Available = defaultAvailable
	}
	if c.Enabled == "" {
		c.Enabled = defaultEnabled
	}

	// Make enabled and available absolute paths
	if !filepath.IsAbs(c.Available) && c.Path != "" {
		c.Available = filepath.Clean(filepath.Join(c.Path, c.Available))
	}
	if !filepath.IsAbs(c.Enabled) && c.Path != "" {
		c.Enabled = filepath.Clean(filepath.Join(c.Path, c.Enabled))
	}

	// Return configuration
	return NewWithConfig(os.DirFS("/").(fs.StatFS), c)
}
