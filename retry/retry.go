// Package retry runs a function with limited retries.
package retry

import (
	"context"
	"errors"
	"time"
)

// Func is the work function. Return a non-nil error to retry (unless attempts exhausted).
type Func func(ctx context.Context) error

// Options configures retry behavior.
type Options struct {
	// Attempts is total tries including the first. Default 3.
	Attempts int
	// Delay is the wait between attempts. Default 0.
	Delay time.Duration
	// Backoff multiplies Delay after each failed attempt. Default 1 (no growth).
	Backoff float64
	// RetryIf decides whether err should be retried. Nil means always retry.
	RetryIf func(error) bool
}

// Do executes fn with the given options until success or attempts are exhausted.
func Do(ctx context.Context, fn Func, opts Options) error {
	if fn == nil {
		return errors.New("retry: nil func")
	}
	if opts.Attempts <= 0 {
		opts.Attempts = 3
	}
	if opts.Backoff <= 0 {
		opts.Backoff = 1
	}
	delay := opts.Delay
	var last error
	for i := 0; i < opts.Attempts; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		last = fn(ctx)
		if last == nil {
			return nil
		}
		if opts.RetryIf != nil && !opts.RetryIf(last) {
			return last
		}
		if i == opts.Attempts-1 {
			break
		}
		if delay > 0 {
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
			delay = time.Duration(float64(delay) * opts.Backoff)
		}
	}
	return last
}
