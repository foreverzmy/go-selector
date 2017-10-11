package selector

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestHasKey(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
	}
	assert.True(HasKey("foo").Matches(valid))
	assert.False(HasKey("zoo").Matches(valid))
}
