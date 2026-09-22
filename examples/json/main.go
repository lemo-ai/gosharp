package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/json"
)

type Order struct {
	ID    int64    `json:"id"`
	Items []string `json:"items"`
}

func main() {
	in := Order{ID: 1001, Items: []string{"a", "b"}}
	_ = json.Pretouch(in)

	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	fmt.Println("Marshal:", string(b))

	var out Order
	if err := json.Unmarshal(b, &out); err != nil {
		panic(err)
	}
	fmt.Printf("Unmarshal: %+v\n", out)
	fmt.Println("Valid:", json.Valid(b))
}
