package gtime_test

import (
	"testing"
	"time"

	"github.com/lemo-ai/gosharp/gtime"
	"github.com/stretchr/testify/assert"
)

func TestSetTimeZoneDoesNotMutateLocal(t *testing.T) {
	defer gtime.ResetTimeZone()
	before := time.Local
	loc, err := gtime.SetTimeZone("Asia/Shanghai")
	assert.NoError(t, err)
	assert.NotNil(t, loc)
	assert.Equal(t, before, time.Local)
	assert.Equal(t, loc, gtime.Location())
}

func TestEqualNilSafe(t *testing.T) {
	var a, b *gtime.Time
	assert.True(t, a.Equal(b))
	assert.False(t, a.Equal(gtime.Now()))
	assert.False(t, gtime.Now().Before(nil))
	assert.Equal(t, time.Duration(0), gtime.Now().Sub(nil))
}

func TestNewNilTimePointer(t *testing.T) {
	var tp *time.Time
	tt := gtime.New(tp)
	assert.NotNil(t, tt)
	assert.True(t, tt.IsZero())
}

func TestStrToTime(t *testing.T) {
	tt, err := gtime.StrToTime("2024-01-02 15:04:05")
	assert.NoError(t, err)
	assert.Equal(t, 2024, tt.Year())
	assert.Equal(t, 1, int(tt.Month()))
	assert.Equal(t, 2, tt.Day())
}

func TestUnmarshalJSONNull(t *testing.T) {
	var tt gtime.Time
	assert.NoError(t, tt.UnmarshalJSON([]byte("null")))
	assert.True(t, tt.IsZero())
}
