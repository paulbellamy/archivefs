package main

import "text/template"

var tmpl = template.Must(template.New("t").Parse(`package {{.Package}}

import "net/http"
import "github.com/paulbellamy/archivefs"

var {{.ExportVar}} http.FileSystem = archivefs.Must(archivefs.FromString(archivefs.Tar, {{.Content | printf "%q"}}))`))
