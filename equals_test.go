package selector

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "bar",
	}
	assert.True(Equals{Key: "foo", Value: "far"}.Matches(valid))
	assert.False(Equals{Key: "zoo", Value: "buzz"}.Matches(valid))
	assert.False(Equals{Key: "foo", Value: "bar"}.Matches(valid))
}
