package archivefs

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Ensure that the filesystem is loaded. If loading the filesystem fails, this
// will panic.
func Must(fs http.FileSystem, err error) http.FileSystem {
	if err != nil {
		panic(err)
	}
	return fs
}

// Load the current executable as an archive. If you've appended your zip file
// of assets to your executable, this is how you should retrieve them.
func CurrentExecutable(f Format) (http.FileSystem, error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}
	return FromFile(f, path)
}

// Load an archive from a file path.
func FromFile(f Format, path string) (http.FileSystem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return FromReader(f, file)
}

// Load an archive from a string of binary data.
func FromString(f Format, content string) (http.FileSystem, error) {
	return FromReader(f, strings.NewReader(content))
}

// Load an archive from an io.Reader.
func FromReader(f Format, r io.Reader) (http.FileSystem, error) {
	return f.FromReader(r)
}
