package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/algo"
)

func main() {
	fmt.Println("GCD(54,24):", algo.GCD(54, 24))
	fmt.Println("LCM(12,18):", algo.LCM(12, 18))
	fmt.Println("IsPrime(97):", algo.IsPrime(97))
	fmt.Println("PrimesBelow(20):", algo.PrimesBelow(20))
	fmt.Println("Fibonacci(10):", algo.Fibonacci(10))

	a := []int{5, 1, 4, 2, 3}
	algo.Sort(a)
	fmt.Println("Sort:", a)
	fmt.Println("BinarySearch 4:", algo.BinarySearch(a, 4))

	fmt.Println("Levenshtein:", algo.Levenshtein("kitten", "sitting"))
	fmt.Println("KMPIndex:", algo.KMPIndex("ababcabcabababd", "ababd"))

	uf := algo.NewUnionFind(4)
	uf.Union(0, 1)
	uf.Union(2, 3)
	fmt.Println("UF Connected(0,1):", uf.Connected(0, 1), "Count:", uf.Count())

	cache := algo.NewLRU[string, int](2)
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Get("a")
	cache.Put("c", 3)
	if _, ok := cache.Get("b"); !ok {
		fmt.Println("LRU evicted b as expected")
	}
}
