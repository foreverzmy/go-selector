package selector

import (
	"fmt"
	"strings"
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

	assert.NotNil(CheckKey(""), "should error on empty keys")

	assert.NotNil(CheckKey("/foo"), "should error on empty dns prefixes")
	superLongDNSPrefixed := fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen), strings.Repeat("a", MaxKeyLen))
	assert.Nil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen+1), strings.Repeat("a", MaxKeyLen))
	assert.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen+1), strings.Repeat("a", MaxKeyLen+1))
	assert.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
	superLongDNSPrefixed = fmt.Sprintf("%s/%s", strings.Repeat("a", MaxDNSPrefixLen), strings.Repeat("a", MaxKeyLen+1))
	assert.NotNil(CheckKey(superLongDNSPrefixed), len(superLongDNSPrefixed))
}

func TestCheckValue(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckValue(""), "should not error on empty values")
	assert.Nil(CheckValue("foo"))
	assert.Nil(CheckValue("bar_baz"))
	assert.NotNil(CheckValue("_bar_baz"))
	assert.NotNil(CheckValue("bar_baz_"))
	assert.NotNil(CheckValue("_bar_baz_"))
}

func TestIsAlpha(t *testing.T) {
	assert := assert.New(t)

	assert.True(isAlpha('A'))
	assert.True(isAlpha('a'))
	assert.True(isAlpha('Z'))
	assert.True(isAlpha('z'))
	assert.True(isAlpha('0'))
	assert.True(isAlpha('9'))
	assert.True(isAlpha('함'))
	assert.True(isAlpha('é'))
	assert.False(isAlpha('-'))
	assert.False(isAlpha('/'))
	assert.False(isAlpha('~'))
}
