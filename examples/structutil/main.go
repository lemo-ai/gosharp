package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/structutil"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fields, err := structutil.TagFields(User{}, []string{"json"})
	if err != nil {
		panic(err)
	}
	for _, f := range fields {
		fmt.Printf("field=%s tag=%s\n", f.Name(), f.TagValue)
	}

	m, err := structutil.TagMapName(User{}, []string{"json"})
	if err != nil {
		panic(err)
	}
	fmt.Println("TagMapName:", m)
}
