package embedfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Executable is used for testing purposes to mock the os.Executable function.
var Executable = os.Executable

// Load returns a live filesystem if called from `go run`, otherwise it returns
// an embedded filesystem.
func Load(embeds fs.FS, dir string) (fs.FS, error) {
	exe, err := Executable()
	if err != nil {
		return nil, err
	}
	// If called from `go run`, use the filesystem
	if strings.Contains(exe, string(filepath.Separator)+"go-build") {
		// Ensure that the directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return nil, err
		}
		// Create a file system rooted at the specified directory
		return os.DirFS(dir), nil
	}
	// Otherwise, use the embedded filesystem, checking if the directory exists
	// first
	if _, err := fs.Stat(embeds, dir); err != nil {
		return nil, err
	}
	// Create a sub filesystem rooted at the specified directory
	fsys, err := fs.Sub(embeds, dir)
	if err != nil {
		return nil, err
	}
	return fsys, nil
}
