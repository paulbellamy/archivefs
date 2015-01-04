# tarfs

Experimental directory to Go net/http.FileSystem generator. Point it at a directory, and it outputs a .go file into your package with the contents of that directory in a net/http.FileSystem var ready to be served.

## Usage

Fetch the library:

```bash
go get github.com/paulbellamy/tarfs
go get github.com/paulbellamy/tarfs/cmd/tarfs
```

In your code, add the invocation:

```golang
//go:generate tarfs -o assets.go -package main -export-var FileSystem assets
```

Then to pack up your assets run:

```bash
go generate
```

## Why

One of the nicest things about Go is single-binary deploys. These are pretty easy until you start adding assets into your project, like html files, images, etc. This lets you bundle your assets into your single binary. The caveat is that when your server boots all your assets will be copied into memory, so it's not a great fit for large assets.

## Todo

* Use a gzipped data format instead of tar to a string.
