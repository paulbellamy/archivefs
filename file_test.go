package archivefs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_File_Stat(t *testing.T) {
	f, err := testFileSystem.Open("subdir/file.txt")
	if !assert.NoError(t, err) {
		return
	}

	info, err := f.Stat()
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "file.txt", info.Name())
	assert.Equal(t, int64(9), info.Size())
	assert.Equal(t, os.FileMode(0x01a4), info.Mode())
	assert.Equal(t, time.Unix(1420380959, 0).Unix(), info.ModTime().Unix())
	assert.False(t, info.IsDir())
	assert.Nil(t, info.Sys())
}
