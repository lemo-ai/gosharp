package regex_test

import (
	"testing"

	"github.com/lemo-ai/gosharp/regex"
	"github.com/stretchr/testify/assert"
)

func TestQuoteAndMatch(t *testing.T) {
	assert.Equal(t, `\[foo\]`, regex.Quote(`[foo]`))
	assert.True(t, regex.IsMatchString(`^\d+$`, "123"))
	assert.False(t, regex.IsMatchString(`^\d+$`, "12a"))
	assert.NoError(t, regex.Validate(`a+`))
	assert.Error(t, regex.Validate(`(`))
}
