package archivefs

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"time"
)

type TarFile struct {
	dir    *TarDir
	header *tar.Header
	body   []byte
}

func (tf *TarFile) Stat() (os.FileInfo, error) {
	return tf, nil
}

// base name of the file
func (tf *TarFile) Name() string {
	return filepath.Base(tf.header.Name)
}

// length in bytes for regular files; system-dependent for others
func (tf *TarFile) Size() int64 {
	return int64(len(tf.body))
}

// file mode bits
func (tf *TarFile) Mode() os.FileMode {
	return os.FileMode(tf.header.Mode)
}

// modification time
func (tf *TarFile) ModTime() time.Time {
	return tf.header.ModTime
}

// abbreviation for Mode().IsDir()
func (tf *TarFile) IsDir() bool {
	return false
}

// underlying data source (can return nil)
func (tf *TarFile) Sys() interface{} {
	return nil
}

// Makes a new reader into this file
func (tf *TarFile) NewReader() *TarFileReader {
	return &TarFileReader{
		TarFile: tf,
		reader:  bytes.NewReader(tf.body),
	}
}

type TarFileReader struct {
	*TarFile
	reader io.ReadSeeker
}

func (tfr *TarFileReader) Read(p []byte) (n int, err error) {
	return tfr.reader.Read(p)
}

func (tfr *TarFileReader) Close() error {
	return nil
}

func (tfr *TarFileReader) Readdir(count int) ([]os.FileInfo, error) {
	return tfr.TarFile.dir.Readdir(count)
}

func (tfr *TarFileReader) Seek(offset int64, whence int) (int64, error) {
	return tfr.reader.Seek(offset, whence)
}

func (tfr *TarFileReader) Stat() (os.FileInfo, error) {
	return tfr.TarFile.Stat()
}
