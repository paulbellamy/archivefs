# archivefs [![Build Status](https://travis-ci.org/paulbellamy/archivefs.svg)]

Directory to Go net/http.FileSystem compiler. Point it at a directory, and it outputs a .go file with the contents of that directory in a net/http.FileSystem var ready to be served. Can also be used as a library if you'd rather roll your own.

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
