package selector

import (
	"testing"

	"github.com/blendlabs/go-assert"
)

func TestParseEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "bar",
		"moo": "lar",
	}
	invalid := Labels{
		"zoo": "mar",
		"moo": "lar",
	}

	selector, err := Parse("foo == bar")
	assert.Nil(err)
	assert.True(selector.Matches(valid))
	assert.False(selector.Matches(invalid))
}

func TestParseNotEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "lar",
	}
	invalidPresent := Labels{
		"foo": "bar",
		"moo": "lar",
	}
	invalidMissing := Labels{
		"zoo": "mar",
		"moo": "lar",
	}

	selector, err := Parse("foo != bar")
	assert.Nil(err)
	assert.True(selector.Matches(valid))
	assert.True(selector.Matches(invalidMissing))
	assert.False(selector.Matches(invalidPresent))
}

func TestParseIn(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "lar",
	}
	valid2 := Labels{
		"foo": "bar",
		"moo": "lar",
	}
	invalid := Labels{
		"foo": "mar",
		"moo": "lar",
	}
	invalidMissing := Labels{
		"zoo": "mar",
		"moo": "lar",
	}

	selector, err := Parse("foo in (bar,far)")
	assert.Nil(err)
	assert.True(selector.Matches(valid), selector.String())
	assert.True(selector.Matches(valid2))
	assert.True(selector.Matches(invalidMissing))
	assert.False(selector.Matches(invalid), selector.String())
}

func TestParseGroup(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"zoo":   "mar",
		"moo":   "lar",
		"thing": "map",
	}
	invalid := Labels{
		"zoo":   "mar",
		"moo":   "something",
		"thing": "map",
	}
	selector, err := Parse("zoo=mar, moo=lar, thing")
	assert.Nil(err)
	assert.True(selector.Matches(valid))
	assert.False(selector.Matches(invalid))
}
