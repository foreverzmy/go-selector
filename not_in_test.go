package selector

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestNotIn(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "mar",
		"moo": "lar",
	}
	invalid := Labels{
		"foo": "far",
		"moo": "lar",
	}
	missing := Labels{
		"loo": "mar",
		"moo": "lar",
	}

	selector := NotIn{Key: "foo", Values: []string{"bar", "far"}}
	assert.True(selector.Matches(valid))
	assert.True(selector.Matches(missing))
	assert.False(selector.Matches(invalid))
	assert.Equal("foo notin (bar, far)", selector.String())
}
