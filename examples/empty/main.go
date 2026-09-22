package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/empty"
)

func main() {
	fmt.Println("nil:", empty.IsEmpty(nil))
	fmt.Println(`"":`, empty.IsEmpty(""))
	fmt.Println("0:", empty.IsEmpty(0))
	fmt.Println("[]int{}:", empty.IsEmpty([]int{}))
	fmt.Println(`"hi":`, empty.IsEmpty("hi"))

	var p *int
	fmt.Println("typed nil pointer:", empty.IsNil(p), empty.IsEmpty(p))
}
