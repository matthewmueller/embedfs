package embedfs_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/matryer/is"
	"github.com/matthewmueller/embedfs"
)

func TestLoadLiveFS(t *testing.T) {
	is := is.New(t)

	// Save and restore original Executable
	origExecutable := embedfs.Executable
	defer func() { embedfs.Executable = origExecutable }()

	// Simulate go run by returning a path containing "go-build"
	embedfs.Executable = func() (string, error) {
		return filepath.Join(os.TempDir(), "go-build", "fake-exe"), nil
	}

	// Use a temp dir as the live FS root
	tmpDir := t.TempDir()
	fsys, err := embedfs.Load(nil, tmpDir)
	is.NoErr(err)
	is.True(fsys != nil)
}

func TestLoadLiveFSDirNotExist(t *testing.T) {
	is := is.New(t)

	origExecutable := embedfs.Executable
	defer func() { embedfs.Executable = origExecutable }()

	embedfs.Executable = func() (string, error) {
		return filepath.Join(os.TempDir(), "go-build", "fake-exe"), nil
	}

	// Use a non-existent directory
	nonExistentDir := filepath.Join(os.TempDir(), "definitely-does-not-exist")
	fsys, err := embedfs.Load(nil, nonExistentDir)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.True(fsys == nil)
}

func TestLoadEmbeddedFS(t *testing.T) {
	is := is.New(t)

	// Save and restore original Executable
	origExecutable := embedfs.Executable
	defer func() { embedfs.Executable = origExecutable }()

	// Simulate normal binary (not go run)
	embedfs.Executable = func() (string, error) {
		return "/usr/local/bin/myapp", nil
	}

	// Create a fake embedded FS
	embeds := fstest.MapFS{
		"foo/bar.txt": &fstest.MapFile{Data: []byte("hello")},
		"foo/baz.txt": &fstest.MapFile{Data: []byte("world")},
	}

	fsys, err := embedfs.Load(embeds, "foo")
	is.NoErr(err)
	is.True(fsys != nil)

	// Should be able to read files from the sub FS
	data, err := fs.ReadFile(fsys, "bar.txt")
	is.NoErr(err)
	is.Equal(string(data), "hello")
}

func TestLoadExecutableError(t *testing.T) {
	is := is.New(t)

	origExecutable := embedfs.Executable
	defer func() { embedfs.Executable = origExecutable }()

	embedfs.Executable = func() (string, error) {
		return "", errors.New("fail")
	}

	fsys, err := embedfs.Load(nil, "foo")
	is.True(err != nil)
	is.True(fsys == nil)
}

func TestLoadEmbeddedFSSubError(t *testing.T) {
	is := is.New(t)

	origExecutable := embedfs.Executable
	defer func() { embedfs.Executable = origExecutable }()

	embedfs.Executable = func() (string, error) {
		return "/usr/local/bin/myapp", nil
	}

	embeds := fstest.MapFS{
		"foo/bar.txt": &fstest.MapFile{Data: []byte("hello")},
	}

	// Try to load a non-existent subdir
	fsys, err := embedfs.Load(embeds, "notfound")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.True(fsys == nil)
}
