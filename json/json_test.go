package json_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/json"
	"github.com/stretchr/testify/assert"
)

func TestMarshalUnmarshal(t *testing.T) {
	type S struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	in := S{A: 1, B: "x"}
	b, err := json.Marshal(in)
	assert.NoError(t, err)
	assert.True(t, json.Valid(b))

	var out S
	assert.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)

	indented, err := json.MarshalIndent(in, "", "  ")
	assert.NoError(t, err)
	assert.Contains(t, string(indented), "\n")
}
