package selector

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestCheckKey(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckKey("foo"))
	assert.Nil(CheckKey("bar/foo"))
	assert.Nil(CheckKey("bar.io/foo"))
	assert.NotNil(CheckKey("_foo"))
	assert.NotNil(CheckKey("-foo"))
	assert.NotNil(CheckKey("foo-"))
	assert.NotNil(CheckKey("foo_"))
	assert.NotNil(CheckKey("bar/foo/baz"))
}

func TestCheckValue(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckValue("foo"))
	assert.Nil(CheckValue("bar_baz"))
	assert.NotNil(CheckValue("_bar_baz"))
	assert.NotNil(CheckValue("bar_baz_"))
	assert.NotNil(CheckValue("_bar_baz_"))
}
