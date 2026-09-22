package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/hashx"
)

func main() {
	s := "gosharp"
	fmt.Println("MD5:", hashx.MD5String(s))
	fmt.Println("SHA1:", hashx.SHA1([]byte(s)))
	fmt.Println("SHA256:", hashx.SHA256String(s))
	fmt.Println("FNV64a:", hashx.FNV64a([]byte(s)))
}
