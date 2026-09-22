package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/lemo-ai/gosharp/safego"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	safego.GoWithHandler(func() {
		defer wg.Done()
		panic("boom")
	}, func(recovered any, stack []byte) {
		fmt.Println("recovered:", recovered)
		fmt.Println("stack bytes:", len(stack))
	})
	wg.Wait()

	wg.Add(1)
	safego.Go(func() {
		defer wg.Done()
		fmt.Println("normal goroutine ok")
	})
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
}
