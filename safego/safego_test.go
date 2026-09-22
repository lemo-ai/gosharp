package safego_test

import (
	"sync"
	"testing"
	"time"

	"github.com/lemo-ai/gosharp/safego"
	"github.com/stretchr/testify/assert"
)

func TestGoRecoversPanic(t *testing.T) {
	var (
		mu   sync.Mutex
		got  any
		done = make(chan struct{})
	)
	safego.GoWithHandler(func() {
		panic("boom")
	}, func(recovered any, _ []byte) {
		mu.Lock()
		got = recovered
		mu.Unlock()
		close(done)
	})
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, "boom", got)
}
