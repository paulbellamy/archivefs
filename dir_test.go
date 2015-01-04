package tarfs

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OpenAndReadFile(t *testing.T) {
	fs, err := New(testData)
	assert.NoError(t, err)

	f, err := fs.Open("readme.txt")
	if !assert.NoError(t, err) {
		return
	}

	body, err := ioutil.ReadAll(f)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "This archive contains some text files.", string(body))
}

func Test_NotFound(t *testing.T) {
	fs, err := New(testData)
	if !assert.NoError(t, err) {
		return
	}

	f, err := fs.Open("not_found.txt")
	assert.EqualError(t, err, "open not_found.txt: file does not exist")
	assert.Nil(t, f)
}

func Test_EscapingSandbox(t *testing.T) {
	fs, err := New(testData)
	if !assert.NoError(t, err) {
		return
	}

	f, err := fs.Open("../../readme.txt")
	assert.EqualError(t, err, "open ../../readme.txt: file does not exist")
	assert.Nil(t, f)
}

func Test_HardRoot(t *testing.T) {
	fs, err := New(testData)
	assert.NoError(t, err)

	f, err := fs.Open("/readme.txt")
	if !assert.NoError(t, err) {
		return
	}

	body, err := ioutil.ReadAll(f)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "This archive contains some text files.", string(body))
}

func Test_HandlesDirs(t *testing.T) {
	fs, err := New(testData)
	if !assert.NoError(t, err) {
		return
	}

	f, err := fs.Open("subdir/file.txt")
	if !assert.NoError(t, err) {
		return
	}

	body, err := ioutil.ReadAll(f)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "sub-file", string(body))
}