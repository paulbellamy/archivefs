package archivefs

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

var Zip = &fnFormat{
	loadFn: decodeZip,
}

func decodeZip(r io.Reader) (http.FileSystem, error) {
	// Open the zip archive for reading.
	// We have to load it all into memory for this step, because zip needs to
	// know how long the reader is.
	zipData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	archive, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	root := &Dir{
		files: map[string]*File{},
	}

	// Iterate through the files in the archive.
	for _, file := range archive.File {
		root.files[file.FileHeader.Name] = &File{
			FileInfo:    file.FileHeader.FileInfo(),
			Dir:         root,
			NewReaderFn: newZipFileReader(file),
		}
	}

	return root, nil
}

func newZipFileReader(zipFile *zip.File) func(*File) (http.File, error) {
	return func(file *File) (http.File, error) {
		readCloser, err := zipFile.Open()
		if err != nil {
			return nil, err
		}

		// Have to add a io.Seeker to readCloser.
		readSeekCloser := NewMemoizingReadCloser(readCloser)

		return &FileReader{
			File:   file,
			reader: readSeekCloser,
			closer: readSeekCloser,
		}, nil
	}
}
