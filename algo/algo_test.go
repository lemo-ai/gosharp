package algo_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/algo"
	"github.com/stretchr/testify/assert"
)

func TestMath(t *testing.T) {
	assert.Equal(t, 3, algo.Min(3, 5))
	assert.Equal(t, 5, algo.Max(3, 5))
	assert.Equal(t, 5, algo.Clamp(10, 1, 5))
	assert.Equal(t, int64(6), algo.GCD(54, 24))
	assert.Equal(t, int64(36), algo.LCM(12, 18))
	assert.Equal(t, int64(8), algo.PowInt(2, 3))
	assert.True(t, algo.IsPrime(97))
	assert.False(t, algo.IsPrime(100))
	assert.Equal(t, []int{2, 3, 5, 7}, algo.PrimesBelow(10))
	assert.Equal(t, int64(55), algo.Fibonacci(10))
}

func TestSearchSort(t *testing.T) {
	a := []int{5, 1, 4, 2, 3}
	algo.Sort(a)
	assert.True(t, algo.IsSorted(a))
	assert.Equal(t, 2, algo.BinarySearch(a, 3))
	assert.Equal(t, -1, algo.BinarySearch(a, 9))
	assert.Equal(t, 2, algo.LowerBound(a, 3))
	assert.Equal(t, 3, algo.UpperBound(a, 3))
}

func TestStringAlgo(t *testing.T) {
	assert.Equal(t, 3, algo.Levenshtein("kitten", "sitting"))
	assert.Equal(t, 4, algo.LCSLength("ABCBDAB", "BDCABA"))
	assert.Equal(t, 2, algo.KMPIndex("hello world", "llo"))
	assert.Equal(t, -1, algo.KMPIndex("hello", "xyz"))
}

func TestUnionFind(t *testing.T) {
	uf := algo.NewUnionFind(5)
	assert.True(t, uf.Union(0, 1))
	assert.True(t, uf.Union(1, 2))
	assert.False(t, uf.Union(0, 2))
	assert.True(t, uf.Connected(0, 2))
	assert.False(t, uf.Connected(0, 3))
	assert.Equal(t, 3, uf.Count())
}

func TestLRU(t *testing.T) {
	c := algo.NewLRU[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	v, ok := c.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	c.Put("c", 3) // evicts b
	_, ok = c.Get("b")
	assert.False(t, ok)
	assert.Equal(t, 2, c.Len())
}
