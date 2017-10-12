package selector

import (
	"testing"

	"github.com/blendlabs/go-assert"
)

func TestParseInvalid(t *testing.T) {
	assert := assert.New(t)

	_, err := Parse("x==a==b")
	assert.NotNil(err)

	_, err = Parse("!x==b")
	assert.NotNil(err)

	_, err = Parse("x<b")
	assert.NotNil(err)
}

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
	invalid2 := Labels{
		"zoo":    "mar",
		"moo":    "lar",
		"!thing": "map",
	}
	selector, err := Parse("zoo=mar, moo=lar, thing")
	assert.Nil(err)
	assert.True(selector.Matches(valid))
	assert.False(selector.Matches(invalid))
	assert.False(selector.Matches(invalid2))

	complicated, err := Parse("zoo in (mar,lar,dar),moo,!thingy")
	assert.Nil(err)
	assert.NotNil(complicated)
	assert.True(complicated.Matches(valid))
}

func TestParseGroupComplicated(t *testing.T) {
	assert := assert.New(t)
	valid := Labels{
		"zoo":   "mar",
		"moo":   "lar",
		"thing": "map",
	}
	complicated, err := Parse("zoo in (mar,lar,dar),moo,thing == map,!thingy")
	assert.Nil(err)
	assert.NotNil(complicated)
	assert.True(complicated.Matches(valid))
}

func TestParseDocsExample(t *testing.T) {
	assert := assert.New(t)
	selector, err := Parse("x in (foo,,baz),y,z notin ()")
	assert.Nil(err)
	assert.NotNil(selector)
}

func TestParseValidate(t *testing.T) {
	assert := assert.New(t)

	_, err := Parse("zoo=bar")
	assert.Nil(err)

	_, err = Parse("_zoo=bar")
	assert.NotNil(err)

	_, err = Parse("_zoo=_bar")
	assert.NotNil(err)

	_, err = Parse("zoo=bar,foo=_mar")
	assert.NotNil(err)
}

func TestParseMultiByte(t *testing.T) {
	assert := assert.New(t)

	selector, err := Parse("함=수,목=록") // number=number, number=rock
	assert.Nil(err)
	assert.NotNil(selector)

	typed, isTyped := selector.(And)
	assert.True(isTyped)
	assert.Len(typed, 2)
}

func BenchmarkParse(b *testing.B) {
	valid := Labels{
		"zoo":   "mar",
		"moo":   "lar",
		"thing": "map",
	}

	for i := 0; i < b.N; i++ {
		selector, err := Parse("zoo in (mar,lar,dar),moo,!thingy")
		if err != nil {
			b.Fail()
		}
		if !selector.Matches(valid) {
			b.Fail()
		}
	}
}
