package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/convert"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fmt.Println("Int:", convert.Int("42"))
	fmt.Println("Float64:", convert.Float64("3.14"))
	fmt.Println("Bool:", convert.Bool("true"))
	fmt.Println("Strings:", convert.Strings([]int{1, 2, 3}))

	var u User
	_ = convert.Struct(map[string]any{"name": "alice", "AGE": "18"}, &u)
	fmt.Printf("Struct: %+v\n", u)
}
