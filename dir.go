package tarfs

import (
	"archive/tar"

	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type TarDir struct {
	dir    *TarDir // parent Dir
	header *tar.Header
	dirs   map[string]*TarDir
	files  map[string]*TarFile
}

func cleanPath(path string) string {
	clean := filepath.Clean(path)
	if filepath.IsAbs(clean) {
		return clean[1:]
	} else {
		return clean
	}
}

func (td *TarDir) Open(name string) (http.File, error) {
	if f, ok := td.files[cleanPath(name)]; ok {
		return f.NewReader(), nil
	}
	return nil, &os.PathError{
		Op:   "open",
		Path: name,
		Err:  os.ErrNotExist,
	}
}

// base name of the dir
func (td *TarDir) Name() string {
	return filepath.Base(td.header.Name)
}

// TODO: Subsequent calls to this should yield further FileInfos, as per os.File
func (fs *TarDir) Readdir(count int) ([]os.FileInfo, error) {
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

func (tfr *TarDir) Stat() (os.FileInfo, error) {
	return &TarDirInfo{TarDir: tfr}, nil
}

type TarDirInfo struct {
	*TarDir
}

// base name of the dir
func (tdi *TarDirInfo) Name() string {
	return tdi.TarDir.Name()
}

// length in bytes for regular files; system-dependent for others
func (tdi *TarDirInfo) Size() int64 {
	return 0
}

// file mode bits
// TODO: return actual permissions
func (tdi *TarDirInfo) Mode() os.FileMode {
	return os.FileMode(tdi.TarDir.header.Mode)
}

// modification time
func (tdi *TarDirInfo) ModTime() time.Time {
	return tdi.TarDir.header.ModTime
}

// abbreviation for Mode().IsDir()
func (tdi *TarDirInfo) IsDir() bool {
	return true
}

// underlying data source (can return nil)
func (tdi *TarDirInfo) Sys() interface{} {
	return nil
}
