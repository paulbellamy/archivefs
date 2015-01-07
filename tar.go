package archivefs

import (
	"archive/tar"
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
