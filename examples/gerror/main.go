package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/lemo-ai/gosharp/gerror"
)

func main() {
	err := gerror.Wrap(io.EOF, "read config")
	fmt.Println("Error:", err)
	fmt.Println("Is EOF:", errors.Is(err, io.EOF))
	fmt.Println("Cause:", gerror.Cause(err))
	fmt.Printf("Stack:\n%s\n", gerror.Stack(err))
}
