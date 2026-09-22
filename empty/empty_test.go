package empty_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/empty"
	"github.com/stretchr/testify/assert"
)

type stringer struct{ s string }

func (s *stringer) String() string { return s.s }

func TestIsEmpty(t *testing.T) {
	assert.True(t, empty.IsEmpty(nil))
	assert.True(t, empty.IsEmpty(""))
	assert.True(t, empty.IsEmpty(0))
	assert.True(t, empty.IsEmpty(false))
	assert.True(t, empty.IsEmpty([]int{}))
	assert.False(t, empty.IsEmpty("x"))
	assert.False(t, empty.IsEmpty(1))
}

func TestIsEmptyTypedNil(t *testing.T) {
	var s *stringer
	var i interface{} = s
	assert.True(t, empty.IsNil(i))
	assert.True(t, empty.IsEmpty(i))
}

func TestIsNil(t *testing.T) {
	var p *int
	assert.True(t, empty.IsNil(p))
	assert.False(t, empty.IsNil(0))
}
