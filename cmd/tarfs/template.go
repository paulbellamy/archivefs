package main

import "text/template"

var tmpl = template.Must(template.New("t").Parse(`
package {{.Package}}

import "net/http"
import "github.com/paulbellamy/httpdir-compiler/tarfs"

var {{.ExportVar}} http.FileSystem = tarfs.Must(tarfs.New({{.Content | printf "%q"}}))
`))
