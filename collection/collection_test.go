package collection_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/collection"
	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	assert.True(t, collection.Contains([]int{1, 2, 3}, 2))
	assert.Equal(t, 1, collection.Index([]string{"a", "b"}, "b"))
	assert.Equal(t, []int{1, 2, 3}, collection.Unique([]int{1, 2, 2, 3, 1}))
	assert.Equal(t, []int{2, 4}, collection.Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 }))
	assert.Equal(t, []int{2, 4, 6}, collection.Map([]int{1, 2, 3}, func(v int) int { return v * 2 }))
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, collection.Chunk([]int{1, 2, 3, 4, 5}, 2))
	assert.Equal(t, []int{3, 2, 1}, collection.Reverse([]int{1, 2, 3}))
	assert.Equal(t, []int{1, 3}, collection.Diff([]int{1, 2, 3}, []int{2}))
	assert.Equal(t, []int{2, 3}, collection.Intersect([]int{1, 2, 3, 2}, []int{2, 3, 4}))
}
