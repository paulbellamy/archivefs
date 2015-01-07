package archivefs

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Dir struct {
	dir    *Dir // parent Dir
	header os.FileInfo
	dirs   map[string]*Dir
	files  map[string]*File
}

func cleanPath(path string) string {
	clean := filepath.Clean(path)
	if filepath.IsAbs(clean) {
		return clean[1:]
	} else {
		return clean
	}
}

func (td *Dir) Open(name string) (http.File, error) {
	if f, ok := td.files[cleanPath(name)]; ok {
		return f.NewReader()
	}
	return nil, &os.PathError{
		Op:   "open",
		Path: name,
		Err:  os.ErrNotExist,
	}
}

// base name of the dir
func (td *Dir) Name() string {
	return filepath.Base(td.header.Name())
}

// TODO: Subsequent calls to this should yield further FileInfos, as per os.File
func (fs *Dir) Readdir(count int) ([]os.FileInfo, error) {
	var results []os.FileInfo
	var added int
	for _, dir := range fs.dirs {
		stat, err := dir.Stat()
		if err != nil {
			return results, err
		}

		results = append(results, stat)

		if count > 0 && added >= count {
			return results, nil
		}
	}

	for _, file := range fs.files {
		stat, err := file.Stat()
		if err != nil {
			return results, err
		}

		results = append(results, stat)

		if count > 0 && added >= count {
			return results, nil
		}
	}

	eof := io.EOF
	if count <= 0 {
		eof = nil
	}

	return results, eof
}

func (tfr *Dir) Stat() (os.FileInfo, error) {
	return &DirInfo{Dir: tfr}, nil
}

type DirInfo struct {
	*Dir
}

// base name of the dir
func (tdi *DirInfo) Name() string {
	return tdi.Dir.Name()
}

// length in bytes for regular files; system-dependent for others
func (tdi *DirInfo) Size() int64 {
	return 0
}

// file mode bits
// TODO: return actual permissions
func (tdi *DirInfo) Mode() os.FileMode {
	return os.FileMode(tdi.Dir.header.Mode())
}

// modification time
func (tdi *DirInfo) ModTime() time.Time {
	return tdi.Dir.header.ModTime()
}

// abbreviation for Mode().IsDir()
func (tdi *DirInfo) IsDir() bool {
	return true
}

// underlying data source (can return nil)
func (tdi *DirInfo) Sys() interface{} {
	return nil
}
