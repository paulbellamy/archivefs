package tarfs

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Must(fs http.FileSystem, err error) http.FileSystem {
	if err != nil {
		panic(err)
	}
	return fs
}

// Read the tar contents and rebuild the file system
func New(content string) (http.FileSystem, error) {
	// Open the tar archive for reading.
	archive := tar.NewReader(strings.NewReader(content))

	root := &TarDir{
		dirs:  map[string]*TarDir{},
		files: map[string]*TarFile{},
	}

	// Iterate through the files in the archive.
	for {
		header, err := archive.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(archive)
		if err != nil {
			return nil, err
		}
		root.files[header.Name] = &TarFile{
			header: header,
			body:   body,
		}
	}

	return root, nil
}
