package tarfs

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Must(fs http.FileSystem, err error) http.FileSystem {
	if err != nil {
		panic(err)
	}
	return fs
}

// Read the tar contents and rebuild the file system
func New(content string) (http.FileSystem, error) {
	// Open the tar archive for reading.
	archive := tar.NewReader(strings.NewReader(content))

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
		fmt.Printf("Contents of %s:\n", header.Name)
		if _, err := io.Copy(os.Stdout, archive); err != nil {
			return nil, err
		}
		fmt.Println()
	}
	return nil, nil
}

type TarDir struct {
	dir    *TarDir // parent Dir
	header *tar.Header
	dirs   map[string]*TarDir
	files  map[string]*TarFile
}

func (fs *TarDir) Open(name string) (http.File, error) {
	return nil, nil
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
