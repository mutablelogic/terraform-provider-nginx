package nginx_test

import (
	"os"
	"path/filepath"
	"testing"

	// Namespace imports
	. "github.com/mutablelogic/terraform-provider-nginx/pkg/nginx"
)

const (
	NGINX_CONG_PATH = "../../etc/nginx"
)

///////////////////////////////////////////////////////////////////////////////

func Test_Folder_000(t *testing.T) {
	if folder, err := NewFolder(FileAbs(t, NGINX_CONG_PATH), true); err != nil {
		t.Error(err)
	} else if _, err := folder.Enumerate(); err != nil {
		t.Error(err)
	}
}

///////////////////////////////////////////////////////////////////////////////

func FileAbs(t *testing.T, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(cwd, path))
}
