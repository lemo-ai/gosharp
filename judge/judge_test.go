package judge_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/judge"
	"github.com/lemo-ai/gosharp/stringutil"
	"github.com/stretchr/testify/assert"
)

func TestIsNumeric(t *testing.T) {
	assert.True(t, judge.IsNumeric("123"))
	assert.True(t, judge.IsNumeric("-12.3"))
	assert.True(t, judge.IsNumeric("+7"))
	assert.False(t, judge.IsNumeric("-"))
	assert.False(t, judge.IsNumeric("+"))
	assert.False(t, judge.IsNumeric("1.2.3"))
	assert.False(t, judge.IsNumeric(".5"))
	assert.False(t, judge.IsNumeric("5."))
	assert.False(t, judge.IsNumeric(""))
}

func TestStringutilDelegates(t *testing.T) {
	assert.Equal(t, judge.IsNumeric("42"), stringutil.IsNumeric("42"))
	assert.Equal(t, judge.UcFirst("abc"), stringutil.UcFirst("abc"))
	assert.Equal(t, judge.RemoveSymbols("a-b_c"), stringutil.RemoveSymbols("a-b_c"))
	assert.True(t, stringutil.EqualFoldWithoutChars("A-B", "ab"))
}
