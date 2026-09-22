package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/collection"
)

func main() {
	nums := []int{1, 2, 2, 3, 4, 4, 5}
	fmt.Println("Unique:", collection.Unique(nums))
	fmt.Println("Contains 3:", collection.Contains(nums, 3))
	fmt.Println("Filter even:", collection.Filter(nums, func(v int) bool { return v%2 == 0 }))
	fmt.Println("Map *10:", collection.Map(nums, func(v int) int { return v * 10 }))
	fmt.Println("Chunk 2:", collection.Chunk([]int{1, 2, 3, 4, 5}, 2))
	fmt.Println("Diff:", collection.Diff([]int{1, 2, 3}, []int{2}))
	fmt.Println("Intersect:", collection.Intersect([]int{1, 2, 3}, []int{2, 3, 4}))
}
