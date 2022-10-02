package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	// Modules
	multierror "github.com/hashicorp/go-multierror"
	config "github.com/mutablelogic/terraform-provider-nginx/pkg/config"
	context "github.com/mutablelogic/terraform-provider-nginx/pkg/context"
	plugin "github.com/mutablelogic/terraform-provider-nginx/pkg/plugin"
	provider "github.com/mutablelogic/terraform-provider-nginx/pkg/provider"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	//. "github.com/mutablelogic/terraform-provider-nginx"
	//. "github.com/mutablelogic/terraform-provider-nginx/plugin"
)

var (
	flagAddr    = flag.String("addr", "", "Address to listen on")
	flagPlugins = flag.String("plugins", "", "Plugin folder")
)

const (
	defaultPluginPattern = "*.plugin"
	pathSeparator        = string(os.PathSeparator)
)

func main() {
	flag.Parse()

	// PluginPath defaults to same folder as executable
	pluginPath, err := GetPluginPath(*flagPlugins)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create a new provider, load plugins
	plugins, err := plugin.LoadPluginsForPattern(pluginPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create context with the address
	ctx := context.ContextForSignal(os.Interrupt, syscall.SIGTERM)
	if *flagAddr != "" {
		ctx = context.WithAddress(ctx, *flagAddr)
	}

	// Use working directory for relative paths
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Get the resources from JSON files
	var result error

	provider := provider.New()
	fs := os.DirFS(string(os.PathSeparator))
	for _, arg := range flag.Args() {
		// Make absolute path
		if !filepath.IsAbs(arg) {
			arg = filepath.Join(wd, arg)
		}

		// Parse JSON files
		resources, err := config.LoadJSONForPattern(fs, strings.TrimPrefix(arg, pathSeparator))
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}

		// Re-parse to create the configuration
		for _, resource := range resources {
			plugin, exists := plugins[resource.Name]
			if !exists {
				result = multierror.Append(result, ErrNotFound.Withf("Plugin not found: %q", resource.Name))
				continue
			}

			// Create a task plugin from the JSON
			plugin, err := config.ParseJSONResource(fs, resource, plugin)
			if err != nil {
				result = multierror.Append(result, err)
				continue
			}

			// Instantiate the plugin into a task
			task, err := provider.New(ctx, plugin)
			if err != nil {
				result = multierror.Append(result, err)
				continue
			}

			fmt.Printf("task=%v\n", task)
		}
	}

	if result != nil {
		fmt.Fprintln(os.Stderr, result)
		os.Exit(1)
	}

	// Run the provider until done
	fmt.Println("Press CTRL+C to exit")
	if err := provider.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func GetPluginPath(defaultPath string) (string, error) {
	if defaultPath == "" {
		if exec, err := os.Executable(); err != nil {
			return "", err
		} else {
			defaultPath = filepath.Join(filepath.Dir(exec), defaultPluginPattern)
		}
	} else if !filepath.IsAbs(defaultPath) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		defaultPath = filepath.Join(wd, defaultPath)
	}
	return defaultPath, nil
}
