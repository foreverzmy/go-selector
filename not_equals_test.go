package selector

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestNotEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "bar",
	}
	assert.False(NotEquals{Key: "foo", Value: "far"}.Matches(valid))
	assert.True(NotEquals{Key: "zoo", Value: "buzz"}.Matches(valid))
	assert.True(NotEquals{Key: "foo", Value: "bar"}.Matches(valid))
}
