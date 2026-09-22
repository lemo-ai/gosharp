package hashx_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/hashx"
	"github.com/stretchr/testify/assert"
)

func TestHashes(t *testing.T) {
	assert.Equal(t, "5d41402abc4b2a76b9719d911017c592", hashx.MD5String("hello"))
	assert.Equal(t, 64, len(hashx.SHA256String("hello")))
	assert.Equal(t, 40, len(hashx.SHA1([]byte("hello"))))
	assert.NotZero(t, hashx.FNV64a([]byte("hello")))
}
