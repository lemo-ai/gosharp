package convert_test

import (
	"testing"
	"time"

	"github.com/lemo-ai/gosharp/convert"
	"github.com/stretchr/testify/assert"
)

func TestInt64Bases(t *testing.T) {
	assert.Equal(t, int64(10), convert.Int64("010"))
	assert.Equal(t, int64(255), convert.Int64("0xff"))
	assert.Equal(t, int64(8), convert.Int64("0o10"))
	assert.Equal(t, int64(-42), convert.Int64("-42"))
}

func TestUint64Bases(t *testing.T) {
	assert.Equal(t, uint64(10), convert.Uint64("010"))
	assert.Equal(t, uint64(255), convert.Uint64("0xFF"))
	assert.Equal(t, uint64(8), convert.Uint64("0O10"))
}

func TestConvertTimePointer(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	got := convert.Convert(now, "*time.Time")
	ptr, ok := got.(*time.Time)
	assert.True(t, ok)
	assert.True(t, now.Equal(*ptr))
}

func TestStringsByteSlice(t *testing.T) {
	assert.Equal(t, []string{"hello"}, convert.Strings([]byte("hello")))
	assert.Equal(t, []string{"a", "b"}, convert.Strings([]string{"a", "b"}))
}

func TestStructBind(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var u User
	err := convert.Struct(map[string]interface{}{
		"name": "alice",
		"age":  "18",
	}, &u)
	assert.NoError(t, err)
	assert.Equal(t, "alice", u.Name)
	assert.Equal(t, 18, u.Age)
}

func TestBool(t *testing.T) {
	assert.True(t, convert.Bool("true"))
	assert.False(t, convert.Bool("false"))
	assert.False(t, convert.Bool("0"))
	assert.True(t, convert.Bool(1))
}
