package archivefs

import (
	"errors"
	"io"
	"io/ioutil"
)

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

func NewMemoizingReadCloser(r io.ReadCloser) *MemoizingReadCloser {
	return &MemoizingReadCloser{
		ReadCloser: r,
		buf:        make([]byte, 0),
		i:          0,
		prevRune:   -1,
	}
}

type MemoizingReadCloser struct {
	io.ReadCloser
	buf      []byte
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}

func (this *MemoizingReadCloser) Read(b []byte) (n int, err error) {
	bLen := len(b)
	if bLen == 0 {
		return 0, nil
	}

	n, err = this.ensureBufferHasLen(int64(bLen) + this.i)

	if this.i >= int64(len(this.buf)) {
		if err == nil {
			err = io.EOF
		}
		return 0, err
	}
	this.prevRune = -1
	n = copy(b, this.buf[this.i:])
	this.i += int64(n)
	return
}

func (this *MemoizingReadCloser) ensureBufferHasLen(needed int64) (n int, err error) {
	if int64(len(this.buf)) >= needed {
		return
	}

	newBuf := make([]byte, needed-int64(len(this.buf)))
	n, err = this.ReadCloser.Read(newBuf)
	if n > 0 {
		this.buf = append(this.buf, newBuf[:n]...)
	}

	return
}

func (this *MemoizingReadCloser) Seek(offset int64, whence int) (n int64, err error) {
	this.prevRune = -1
	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = int64(this.i) + offset
	case 2:
		_, err = ioutil.ReadAll(this) // realize the lazy reader so we can find the end
		abs = int64(len(this.buf)) + offset
	default:
		return 0, errors.New("archivefs.MemoizingReader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("archivefs.MemoizingReader.Seek: negative position")
	}
	this.i = abs

	if err == io.EOF {
		err = nil
	}
	return abs, err
}

func (this *MemoizingReadCloser) Rewind() {
	this.Seek(0, 0)
}
