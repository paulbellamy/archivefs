package archivefs

import "net/http"

var testFileSystems = []http.FileSystem{testTarFileSystem, testZipFileSystem}
