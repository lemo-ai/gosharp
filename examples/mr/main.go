package main

import (
	"fmt"

	"github.com/lemo-ai/gosharp/mr"
)

func main() {
	v, err := mr.MapReduce(
		func(source chan<- interface{}) {
			for i := 1; i <= 5; i++ {
				source <- i
			}
		},
		func(item interface{}, writer mr.Writer, cancel func(error)) {
			writer.Write(item.(int) * item.(int))
		},
		func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
			sum := 0
			for x := range pipe {
				sum += x.(int)
			}
			writer.Write(sum)
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("sum of squares 1..5 =", v)

	err = mr.Finish(
		func() error { fmt.Println("task A"); return nil },
		func() error { fmt.Println("task B"); return nil },
	)
	fmt.Println("Finish err:", err)
}
