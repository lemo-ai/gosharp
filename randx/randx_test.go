package randx_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/randx"
	"github.com/stretchr/testify/assert"
)

func TestRandx(t *testing.T) {
	n, err := randx.Intn(10)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, n, 0)
	assert.Less(t, n, 10)

	s, err := randx.String(16)
	assert.NoError(t, err)
	assert.Len(t, s, 16)

	h, err := randx.Hex(8)
	assert.NoError(t, err)
	assert.Len(t, h, 16)

	assert.Less(t, randx.FastIntn(5), 5)
	_, err = randx.Uint64()
	assert.NoError(t, err)
}
