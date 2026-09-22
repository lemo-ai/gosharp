package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/regex"
)

func main() {
	fmt.Println("Quote:", regex.Quote(`[foo]`))
	fmt.Println("IsMatchString digits:", regex.IsMatchString(`^\d+$`, "12345"))
	fmt.Println("IsMatchString digits:", regex.IsMatchString(`^\d+$`, "12a"))

	m, err := regex.MatchString(`(\w+)@(\w+)`, "a@b")
	if err != nil {
		panic(err)
	}
	fmt.Println("MatchString:", m)
}
