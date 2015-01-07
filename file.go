package archivefs

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	dir    *Dir
	header os.FileInfo
	body   []byte
}

func (tf *File) Stat() (os.FileInfo, error) {
	return tf, nil
}

// base name of the file
func (tf *File) Name() string {
	return filepath.Base(tf.header.Name())
}

// length in bytes for regular files; system-dependent for others
func (tf *File) Size() int64 {
	return int64(len(tf.body))
}

// file mode bits
func (tf *File) Mode() os.FileMode {
	return os.FileMode(tf.header.Mode())
}

// modification time
func (tf *File) ModTime() time.Time {
	return tf.header.ModTime()
}

// abbreviation for Mode().IsDir()
func (tf *File) IsDir() bool {
	return false
}

// underlying data source (can return nil)
func (tf *File) Sys() interface{} {
	return nil
}

// Makes a new reader into this file
func (tf *File) NewReader() (http.File, error) {
	return &FileReader{
		File:   tf,
		reader: bytes.NewReader(tf.body),
	}, nil
}

type FileReader struct {
	*File
	reader io.ReadSeeker
}

func (tfr *FileReader) Read(p []byte) (n int, err error) {
	return tfr.reader.Read(p)
}

func (tfr *FileReader) Close() error {
	return nil
}

func (tfr *FileReader) Readdir(count int) ([]os.FileInfo, error) {
	return tfr.File.dir.Readdir(count)
}

func (tfr *FileReader) Seek(offset int64, whence int) (int64, error) {
	return tfr.reader.Seek(offset, whence)
}

func (tfr *FileReader) Stat() (os.FileInfo, error) {
	return tfr.File.Stat()
}
