package archivefs

import (
	"net/http"
	"os"
	"path/filepath"
)

type File struct {
	os.FileInfo
	Dir         *Dir
	NewReaderFn func(*File) (http.File, error)
}

func (tf *File) Stat() (os.FileInfo, error) {
	return tf.FileInfo, nil
}

// base name of the file
func (tf *File) Name() string {
	return filepath.Base(tf.FileInfo.Name())
}

// Makes a new reader into this file
func (tf *File) NewReader() (http.File, error) {
	return tf.NewReaderFn(tf)
}
