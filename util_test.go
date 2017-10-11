package selector

import (
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestCheckLabel(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckLabel("foo"))
	assert.Nil(CheckLabel("bar/foo"))
	assert.Nil(CheckLabel("bar.io/foo"))
	assert.NotNil(CheckLabel("_foo"))
	assert.NotNil(CheckLabel("-foo"))
	assert.NotNil(CheckLabel("foo-"))
	assert.NotNil(CheckLabel("foo_"))
	assert.NotNil(CheckLabel("bar/foo/baz"))
}

func TestCheckValue(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckValue("foo"))
	assert.Nil(CheckValue("bar_baz"))
	assert.NotNil(CheckValue("_bar_baz"))
	assert.NotNil(CheckValue("bar_baz_"))
	assert.NotNil(CheckValue("_bar_baz_"))
}
