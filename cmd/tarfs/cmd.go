package main

import (
	"archive/tar"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var err error
	flags := flag.NewFlagSet("httpdir-compiler", flag.ExitOnError)
	exportVar := flags.String("export-var", "Dir", "Name of the exported var.")
	pkg := flags.String("package", "", "Name of the generated package. Default is the directory name of the target.")
	target := flags.String("o", "", "Name of a file to write the output to. Default is stdout.")
	flags.Parse(os.Args[1:])

	source := flags.Arg(0)
	if source == "" {
		source, err = os.Getwd()
		checkErr(err)
	}

	output := os.Stdout
	if *target != "" {
		output, err = os.Create(*target)
		checkErr(err)
	}

	if *pkg == "" {
		*pkg = filepath.Base(source)
	}

	if len(*exportVar) < 1 {
		exit(errors.New("export-var cannot be blank"))
	}

	// TODO: Gzip it instead of string
	content, err := tarDir(source)
	checkErr(err)

	checkErr(tmpl.Execute(output, Options{
		Package:   *pkg,
		ExportVar: *exportVar,
		Content:   content,
	}))
}

func checkErr(err error) {
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

type Options struct {
	Package   string
	ExportVar string
	Content   string
}

// Reads the source dir, and tars it into a big string.
// TODO: Read the source dir instead of hardcoded data
func tarDir(source string) (string, error) {
	// Create a buffer to write our archive to.
	buffer := &bytes.Buffer{}

	// Create a new tar archive.
	archive := tar.NewWriter(buffer)

	// Add some files to the archive.
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Println("Adding:", cleanPath(strings.TrimPrefix(path, source)))
		header := &tar.Header{
			Name:    cleanPath(strings.TrimPrefix(path, source)),
			Size:    info.Size(),
			Mode:    int64(info.Mode().Perm()),
			ModTime: info.ModTime(),
		}
		if err := archive.WriteHeader(header); err != nil {
			return err
		}
		if _, err := io.Copy(archive, file); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	// Make sure to check the error on Close.
	if err := archive.Close(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func cleanPath(path string) string {
	return filepath.Clean(strings.TrimPrefix(path, fmt.Sprintf("%c", filepath.Separator)))
}
