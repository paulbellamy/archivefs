package tarfs

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
)

type TarFile struct {
	dir    *TarDir
	header *tar.Header
	body   []byte
}

func (tf *TarFile) Stat() (os.FileInfo, error) {
	return nil, nil
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
