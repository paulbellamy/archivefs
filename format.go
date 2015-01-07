package archivefs

import (
	"io"
	"net/http"
)

type Format interface {
	FromReader(io.Reader) (http.FileSystem, error)
}

type fnFormat struct {
	loadFn func(io.Reader) (http.FileSystem, error)
}

func (f *fnFormat) FromReader(r io.Reader) (http.FileSystem, error) {
	return f.loadFn(r)
}
