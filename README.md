# archivefs ![Build Status](https://travis-ci.org/paulbellamy/archivefs.svg)

Serve Tar and Zip archives as net/http.FileSystems. Let's you compile assets into your static binary for deployment, or zip them and append it to the binary. Still a work in progress.

## Documentation

[http://godoc.org/github.com/paulbellamy/archivefs](http://godoc.org/github.com/paulbellamy/archivefs)

## Usage

Fetch the library:

```bash
# install the library
go get github.com/paulbellamy/archivefs

# install the cli command
go get github.com/paulbellamy/archivefs/cmd/archivefs
```

In your code, add the invocation:

```golang
//go:generate archivefs -o assets.go -package main -export-var Assets assets_directory
```

Then to pack up your assets run:

```bash
go generate
```

To serve up assets:

```golang
package main

import (
  "log"
  "net/http"
)

func main() {
  log.Fatal(http.ListenAndServe(":8080", http.FileServer(Assets)))
}
```

### Usage as a library

archivefs can also be used as a library, without the command. You might want to do this if you already have an archive of your assets, if your archive is compressed, or if your archive is appended to your executable file (instead of compiled into it). This is also the way to use it for Zip archives.

## Why

One of the nicest things about Go is single-binary deploys. These are pretty easy until you start adding assets into your project, like html files, images, etc. This lets you bundle your assets into your single binary. The caveat with tar is that when your server boots all your assets will be copied into memory, so zip is a much better fit for large assets.
