package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/randx"
)

func main() {
	n, err := randx.Intn(100)
	if err != nil {
		panic(err)
	}
	fmt.Println("Intn(100):", n)

	s, err := randx.String(12)
	if err != nil {
		panic(err)
	}
	fmt.Println("String(12):", s)

	h, err := randx.Hex(8)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hex(8):", h)
	fmt.Println("FastIntn(10):", randx.FastIntn(10))
}
