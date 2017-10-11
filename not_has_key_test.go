package selector

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestNotHasKey(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
	}
	assert.False(NotHasKey("foo").Matches(valid))
	assert.True(NotHasKey("zoo").Matches(valid))
}
