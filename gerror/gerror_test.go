package gerror_test

import (
	"errors"
	"io"
	"testing"

	"github.com/lemo-ai/gosharp/gerror"
	"github.com/stretchr/testify/assert"
)

func TestUnwrapAndIs(t *testing.T) {
	err := gerror.Wrap(io.EOF, "read failed")
	assert.True(t, errors.Is(err, io.EOF))
	assert.Equal(t, io.EOF, errors.Unwrap(err))
}

func TestCause(t *testing.T) {
	root := gerror.New("root")
	wrapped := gerror.Wrap(root, "outer")
	assert.Equal(t, "root", gerror.Cause(wrapped).Error())

	only := gerror.New("alone")
	assert.Equal(t, only, gerror.Cause(only))
}

func TestStack(t *testing.T) {
	err := gerror.New("boom")
	assert.NotEmpty(t, gerror.Stack(err))
	assert.Contains(t, err.Error(), "boom")
}

func TestNewEmpty(t *testing.T) {
	assert.Nil(t, gerror.New(""))
}
