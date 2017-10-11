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
