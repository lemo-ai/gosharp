package structutil_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/structutil"
	"github.com/stretchr/testify/assert"
)

type base struct {
	ID int `json:"id"`
}

type user struct {
	base
	Name string `json:"name"`
	Dup  string `json:"id"` // duplicate tag with embedded
}

func TestTagFieldsDedup(t *testing.T) {
	fields, err := structutil.TagFields(user{}, []string{"json"})
	assert.NoError(t, err)
	tags := make(map[string]int)
	for _, f := range fields {
		tags[f.TagValue]++
	}
	assert.Equal(t, 1, tags["id"], "duplicate tag should be filtered")
	assert.Equal(t, 1, tags["name"])
}
