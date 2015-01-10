package archivefs

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReadCloser struct {
	mock.Mock
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	args := m.Mock.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReadCloser) Seek(offset int64, whence int) (int64, error) {
	args := m.Mock.Called(offset)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockReadCloser) Close() error {
	args := m.Mock.Called()
	return args.Error(0)
}

func Test_MemoizingReadCloser(t *testing.T) {
	expected := "hello world"
	r := NewMemoizingReadCloser(ioutil.NopCloser(strings.NewReader(expected)))
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))
}

func Test_MemoizingReader_SeeksBackwards(t *testing.T) {
	expected := "hello world"
	r := NewMemoizingReadCloser(ioutil.NopCloser(strings.NewReader(expected)))
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))

	n, err := r.Seek(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	b, err = ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))
}

func Test_MemoizingReader_Rewinds(t *testing.T) {
	expected := "hello world"
	r := NewMemoizingReadCloser(ioutil.NopCloser(strings.NewReader(expected)))
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))

	r.Rewind()

	b, err = ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))
}

func Test_MemoizingReader_PassesThroughErrors(t *testing.T) {
	m := &mockReadCloser{}
	b := make([]byte, 61)
	m.On("Read", b).Return(0, errors.New("test error"))

	n, err := NewMemoizingReadCloser(m).Read(b)
	assert.Equal(t, 0, n)
	assert.EqualError(t, err, "test error")

	m.AssertExpectations(t)
}

func Test_MemoizingReader_PassesThroughClose(t *testing.T) {
	m := &mockReadCloser{}
	m.On("Close").Return(nil)

	err := NewMemoizingReadCloser(m).Close()
	assert.NoError(t, err)

	m.AssertExpectations(t)
}

func Test_MemoizingReader_PassesThroughCloseErrors(t *testing.T) {
	m := &mockReadCloser{}
	m.On("Close").Return(errors.New("test error"))

	err := NewMemoizingReadCloser(m).Close()
	assert.EqualError(t, err, "test error")

	m.AssertExpectations(t)
}

func Test_MemoizingReader_InvalidWhencesAreIllegal(t *testing.T) {
	n, err := NewMemoizingReadCloser(nil).Seek(0, 4)
	assert.EqualError(t, err, "archivefs.MemoizingReader.Seek: invalid whence")
	assert.Equal(t, 0, n)
}

func Test_MemoizingReader_NegativeSeeksAreIllegal(t *testing.T) {
	n, err := NewMemoizingReadCloser(nil).Seek(-1, 0)
	assert.EqualError(t, err, "archivefs.MemoizingReader.Seek: negative position")
	assert.Equal(t, 0, n)
}

func Test_MemoizingReader_SeeksPastTheEndMoveToTheEndAndRead(t *testing.T) {
	expected := "hello world"
	r := NewMemoizingReadCloser(ioutil.NopCloser(strings.NewReader(expected)))
	n, err := r.Seek(5, 2) // 5 past the end
	assert.Equal(t, len(expected)+5, n)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, "", string(b))
}

func Test_MemoizingReader_ReturnsErrorsWhenSeekingPastEnd(t *testing.T) {
	m := &mockReadCloser{}
	m.On("Read", mock.AnythingOfType("[]uint8")).Return(5, errors.New("test error"))

	r := NewMemoizingReadCloser(m)
	n, err := r.Seek(5, 2) // 5 past the end
	assert.Equal(t, 10, n)
	assert.EqualError(t, err, "test error")

	m.AssertExpectations(t)
}
