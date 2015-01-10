package archivefs

import (
	"io"
	"os"
)

type FileReader struct {
	*File
	reader io.ReadSeeker
	closer io.Closer
}

func (tfr *FileReader) Read(p []byte) (n int, err error) {
	return tfr.reader.Read(p)
}

func (tfr *FileReader) Close() error {
	if tfr.closer != nil {
		return tfr.closer.Close()
	}
	return nil
}

func (tfr *FileReader) Readdir(count int) ([]os.FileInfo, error) {
	return tfr.File.Dir.Readdir(count)
}

func (tfr *FileReader) Seek(offset int64, whence int) (int64, error) {
	return tfr.reader.Seek(offset, whence)
}

func (tfr *FileReader) Stat() (os.FileInfo, error) {
	return tfr.File.Stat()
}
