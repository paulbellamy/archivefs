package archivefs

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func Must(fs http.FileSystem, err error) http.FileSystem {
	if err != nil {
		panic(err)
	}
	return fs
}

func CurrentExecutable(f Format) (http.FileSystem, error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}
	return FromFile(f, path)
}

func FromFile(f Format, path string) (http.FileSystem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return FromReader(f, file)
}

// Read the tar contents and rebuild the file system
func FromString(f Format, content string) (http.FileSystem, error) {
	return FromReader(f, strings.NewReader(content))
}

func FromReader(f Format, r io.Reader) (http.FileSystem, error) {
	return f.FromReader(r)
}
