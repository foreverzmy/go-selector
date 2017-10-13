package selector

import (
	"strings"
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestLexerIsWhitespace(t *testing.T) {
	assert := assert.New(t)

	l := &Lexer{}
	assert.True(l.isWhitespace(' '))
	assert.True(l.isWhitespace('\n'))
	assert.True(l.isWhitespace('\r'))
	assert.True(l.isWhitespace('\t'))

	assert.False(l.isWhitespace('a'))
	assert.False(l.isWhitespace('z'))
	assert.False(l.isWhitespace('A'))
	assert.False(l.isWhitespace('Z'))
	assert.False(l.isWhitespace('1'))
	assert.False(l.isWhitespace('-'))
}

func TestLexerIsAlpha(t *testing.T) {
	assert := assert.New(t)

	l := &Lexer{}
	assert.True(l.isAlpha('a'))
	assert.True(l.isAlpha('z'))
	assert.True(l.isAlpha('A'))
	assert.True(l.isAlpha('Z'))
	assert.True(l.isAlpha('1'))

	assert.False(l.isAlpha('-'))
	assert.False(l.isAlpha(' '))
	assert.False(l.isAlpha('\n'))
	assert.False(l.isAlpha('\r'))
	assert.False(l.isAlpha('\t'))
}

func TestLexerSkipWhitespace(t *testing.T) {
	assert := assert.New(t)

	l := &Lexer{s: "foo    != bar    ", pos: 3}
	assert.Equal(" ", string(l.current()))
	l.skipWhiteSpace()
	assert.Equal(7, l.pos)
	assert.Equal("!", string(l.current()))
	l.pos = 14
	assert.Equal(" ", string(l.current()))
	l.skipWhiteSpace()
	assert.Equal(len(l.s), l.pos)
}

func TestLexerReadWord(t *testing.T) {
	assert := assert.New(t)

	l := &Lexer{s: "foo != bar"}
	assert.Equal("foo", l.readWord())
	assert.Equal(" ", string(l.current()))

	l = &Lexer{s: "foo,"}
	assert.Equal("foo", l.readWord())
	assert.Equal(",", string(l.current()))

	l = &Lexer{s: "foo"}
	assert.Equal("foo", l.readWord())
	assert.True(l.done())
}

func TestLexerReadOp(t *testing.T) {
	assert := assert.New(t)

	l := &Lexer{s: "!= bar"}
	op, err := l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal(" ", string(l.current()))

	l = &Lexer{s: "!=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal("b", string(l.current()))

	l = &Lexer{s: "!=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.Equal("b", string(l.current()))

	l = &Lexer{s: "!="}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("!=", op)
	assert.True(l.done())

	l = &Lexer{s: "= bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("=", op)
	assert.Equal(" ", string(l.current()))

	l = &Lexer{s: "=bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("=", op)
	assert.Equal("b", string(l.current()))

	l = &Lexer{s: "== bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("==", op)
	assert.Equal(" ", string(l.current()))

	l = &Lexer{s: "==bar"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("==", op)
	assert.Equal("b", string(l.current()))

	l = &Lexer{s: "in (foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("in", op)
	assert.Equal(" ", string(l.current()))

	l = &Lexer{s: "in(foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("in", op)
	assert.Equal("(", string(l.current()))

	l = &Lexer{s: "notin (foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("notin", op)
	assert.Equal(" ", string(l.current()))

	l = &Lexer{s: "notin(foo)"}
	op, err = l.readOp()
	assert.Nil(err)
	assert.Equal("notin", op)
	assert.Equal("(", string(l.current()))
}

func TestLexerReadCSV(t *testing.T) {
	assert := assert.New(t)

	l := &Lexer{s: "(bar, baz, biz)"}
	words, err := l.readCSV()
	assert.Nil(err)
	assert.Len(words, 3, strings.Join(words, ","))
	assert.Equal("bar", words[0])
	assert.Equal("baz", words[1])
	assert.Equal("biz", words[2])
	assert.True(l.done())

	l = &Lexer{s: "(bar, buzz, baz"}
	words, err = l.readCSV()
	assert.NotNil(err)

	l = &Lexer{s: "()"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Empty(words)
	assert.True(l.done())

	l = &Lexer{s: "(), thing=after"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Empty(words)
	assert.Equal(",", string(l.current()))

	l = &Lexer{s: "(foo, bar), buzz=light"}
	words, err = l.readCSV()
	assert.Nil(err)
	assert.Len(words, 2)
	assert.Equal("foo", words[0])
	assert.Equal("bar", words[1])
	assert.Equal(",", string(l.current()))

	l = &Lexer{s: "(test, space are bad)"}
	words, err = l.readCSV()
	assert.NotNil(err)
}

func TestLexerHasKey(t *testing.T) {
	assert := assert.New(t)
	l := &Lexer{s: "foo"}
	valid, err := l.Lex()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(HasKey)
	assert.True(isTyped)
	assert.Equal("foo", string(typed))
}

func TestLexerNotHasKey(t *testing.T) {
	assert := assert.New(t)
	l := &Lexer{s: "!foo"}
	valid, err := l.Lex()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(NotHasKey)
	assert.True(isTyped)
	assert.Equal("foo", string(typed))
}

func TestLexerEquals(t *testing.T) {
	assert := assert.New(t)

	l := &Lexer{s: "foo = bar"}
	valid, err := l.Lex()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)

	l = &Lexer{s: "foo=bar"}
	valid, err = l.Lex()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped = valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestLexerDoubleEquals(t *testing.T) {
	assert := assert.New(t)
	l := &Lexer{s: "foo == bar"}
	valid, err := l.Lex()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(Equals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestLexerNotEquals(t *testing.T) {
	assert := assert.New(t)
	l := &Lexer{s: "foo != bar"}
	valid, err := l.Lex()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(NotEquals)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Equal("bar", typed.Value)
}

func TestLexerIn(t *testing.T) {
	assert := assert.New(t)
	l := &Lexer{s: "foo in (bar, baz)"}
	valid, err := l.Lex()
	assert.Nil(err)
	assert.NotNil(valid)
	typed, isTyped := valid.(In)
	assert.True(isTyped)
	assert.Equal("foo", typed.Key)
	assert.Len(typed.Values, 2)
	assert.Equal("bar", typed.Values[0])
	assert.Equal("baz", typed.Values[1])
}

func TestLexerLex(t *testing.T) {
	assert := assert.New(t)
	l := &Lexer{s: ""}
	_, err := l.Lex()
	assert.NotNil(err)
}
