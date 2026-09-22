package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lemo-ai/gosharp/retry"
)

func main() {
	attempts := 0
	err := retry.Do(context.Background(), func(ctx context.Context) error {
		attempts++
		fmt.Println("try", attempts)
		if attempts < 3 {
			return errors.New("temporary")
		}
		return nil
	}, retry.Options{
		Attempts: 5,
		Delay:    20 * time.Millisecond,
		Backoff:  2,
	})
	fmt.Println("done, err=", err, "attempts=", attempts)
}
