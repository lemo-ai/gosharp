package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/judge"
)

func main() {
	fmt.Println("IsNumeric 12.3:", judge.IsNumeric("12.3"))
	fmt.Println("IsNumeric -:", judge.IsNumeric("-"))
	fmt.Println("UcFirst:", judge.UcFirst("gosharp"))
	fmt.Println("RemoveSymbols:", judge.RemoveSymbols("a-b_c.d"))
	fmt.Println("EqualFoldWithoutChars:", judge.EqualFoldWithoutChars("A-B", "ab"))
}
