package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lemo-ai/gosharp/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryEventuallySucceeds(t *testing.T) {
	var n int32
	err := retry.Do(context.Background(), func(ctx context.Context) error {
		if atomic.AddInt32(&n, 1) < 3 {
			return errors.New("tmp")
		}
		return nil
	}, retry.Options{Attempts: 5, Delay: time.Millisecond})
	assert.NoError(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&n))
}

func TestRetryGiveUp(t *testing.T) {
	errDummy := errors.New("no")
	err := retry.Do(context.Background(), func(ctx context.Context) error {
		return errDummy
	}, retry.Options{Attempts: 2, RetryIf: func(error) bool { return false }})
	assert.Equal(t, errDummy, err)
}
