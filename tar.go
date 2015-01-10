package archivefs

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

var Tar = &fnFormat{
	loadFn: decodeTar,
}

func decodeTar(r io.Reader) (http.FileSystem, error) {
	// Open the tar archive for reading.
	archive := tar.NewReader(r)

	root := &Dir{
		dirs:  map[string]*Dir{},
		files: map[string]*File{},
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
		root.files[header.Name] = &File{
			FileInfo:    header.FileInfo(),
			Dir:         root,
			NewReaderFn: newTarFileReader(body),
		}
	}

	return root, nil
}

func newTarFileReader(body []byte) func(*File) (http.File, error) {
	return func(file *File) (http.File, error) {
		return &FileReader{
			File:   file,
			reader: bytes.NewReader(body),
		}, nil
	}
}
