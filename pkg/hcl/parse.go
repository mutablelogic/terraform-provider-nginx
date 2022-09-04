package hcl

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	// Module imports
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	fileExtHCL  = ".hcl"
	fileExtJSON = ".json"
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Parse(filesys fs.FS, path string) (hcl.Body, error) {
	parser := hclparse.NewParser()
	if err := fs.WalkDir(filesys, strings.TrimPrefix(path, string(os.PathSeparator)), func(path string, info fs.DirEntry, err error) error {
		return walkconfig(parser, filesys, path, info, err)
	}); err != nil {
		return nil, err
	}

	// Merge files together
	files := parser.Files()
	body := make([]*hcl.File, 0, len(files))
	for _, file := range files {
		body = append(body, file)
	}

	// Return success
	return hcl.MergeFiles(body), nil
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func walkconfig(parser *hclparse.Parser, filesys fs.FS, path string, d fs.DirEntry, err error) error {
	// Pass down any error
	if err != nil {
		return err
	}
	// Ignore any hidden files
	if strings.HasPrefix(d.Name(), ".") {
		if d.IsDir() {
			return fs.SkipDir
		} else {
			return nil
		}
	}
	// Recurse into folders
	if d.IsDir() {
		return nil
	}
	// Return error if not a regular file
	if !d.Type().IsRegular() {
		return ErrNotImplemented.Withf("%q", d.Name())
	}
	// Deal with filetypes
	switch strings.ToLower(filepath.Ext(d.Name())) {
	case fileExtHCL:
		if data, err := readall(filesys, path); err != nil {
			return err
		} else if _, diags := parser.ParseHCL(data, path); diags.HasErrors() {
			return diags
		}
	case fileExtJSON:
		if data, err := readall(filesys, path); err != nil {
			return err
		} else if _, diags := parser.ParseJSON(data, path); diags.HasErrors() {
			return diags
		}
	default:
		return ErrNotImplemented.Withf("%q", d.Name())
	}
	return nil
}

// Read all bytes from regular file
func readall(filesys fs.FS, path string) ([]byte, error) {
	return fs.ReadFile(filesys, path)
}
