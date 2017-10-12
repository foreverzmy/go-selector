package selector

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/blendlabs/go-exception"
)

const (
	// At is a common rune.
	At = rune('@')
	// Colon is a common rune.
	Colon = rune(':')
	// Dash is a common rune.
	Dash = rune('-')
	// Underscore  is a common rune.
	Underscore = rune('_')
	// Dot is a common rune.
	Dot = rune('.')
	// ForwardSlash is a common rune.
	ForwardSlash = rune('/')
	// BackSlash is a common rune.
	BackSlash = rune('\\')
	// BackTick is a common rune.
	BackTick = rune('`')
	// Bang is a common rune.
	Bang = rune('!')
	// Comma is a common rune.
	Comma = rune(',')
	// OpenBracket is a common rune.
	OpenBracket = rune('[')
	// OpenParens is a common rune.
	OpenParens = rune('(')
	// OpenCurly is a common rune.
	OpenCurly = rune('{')
	// CloseBracket is a common rune.
	CloseBracket = rune(']')
	// CloseParens is a common rune.
	CloseParens = rune(')')
	// Equal is a common rune.
	Equal = rune('=')
	// Space is a common rune.
	Space = rune(' ')
	// Tab is a common rune.
	Tab = rune('\t')
	// Tilde is a common rune.
	Tilde = rune('~')
	// CarriageReturn is a common rune.
	CarriageReturn = rune('\r')
	// NewLine is a common rune.
	NewLine = rune('\n')
)

var (
	// ErrEmptySelector is returned if the selector to be compiled is empty.
	ErrEmptySelector = fmt.Errorf("empty selector")

	// ErrInvalidOperator is returned if the operator is invalid.
	ErrInvalidOperator = fmt.Errorf("invalid operator")

	// ErrInvalidSelector is returned if there is a structural issue with the selector.
	ErrInvalidSelector = fmt.Errorf("invalid selector")

	// ErrKeyEmpty indicates a key is empty.
	ErrKeyEmpty = fmt.Errorf("key empty")

	// ErrKeyTooLong indicates a key is too long.
	ErrKeyTooLong = fmt.Errorf("key too long")

	// ErrKeyDNSPrefixEmpty indicates a key's "dns" prefix is empty.
	ErrKeyDNSPrefixEmpty = fmt.Errorf("key dns prefix empty")

	// ErrKeyDNSPrefixTooLong indicates a key's "dns" prefix is empty.
	ErrKeyDNSPrefixTooLong = fmt.Errorf("key dns prefix too long; must be less than 253 characters")

	// ErrValueTooLong indicates a value is too long.
	ErrValueTooLong = fmt.Errorf("value too long; must be less than 63 characters")

	// ErrKeyInvalidCharacter indicates a key contains characters
	ErrKeyInvalidCharacter = fmt.Errorf(`key contains invalid characters, regex used: ([A-Za-z0-9_-\.])`)

	// MaxDNSPrefixLen is the maximum dns prefix length.
	MaxDNSPrefixLen = 253
	// MaxKeyLen is the maximum key length.
	MaxKeyLen = 63
	// MaxValueLen is the maximum value length.
	MaxValueLen = 63

	// MaxKeyTotalLen is the maximum total key length.
	MaxKeyTotalLen = MaxDNSPrefixLen + MaxKeyLen + 1
)

// CheckKey validates a key.
func CheckKey(key string) (err error) {
	keyLen := len(key)
	if keyLen == 0 {
		err = ErrKeyEmpty
		return
	}
	if keyLen > MaxKeyTotalLen {
		err = ErrKeyTooLong
		return
	}

	var working []rune
	var state int
	var ch rune
	var width int
	for pos := 0; pos < keyLen; pos += width {
		ch, width = utf8.DecodeRuneInString(key[pos:])
		switch state {
		case 0: // collect dns prefix or key
			if ch == ForwardSlash {
				err = checkDNS(string(working))
				if err != nil {
					return
				}
				working = nil
				state = 1
				continue
			}
		}
		working = append(working, ch)
		continue
	}

	if len(working) > MaxKeyLen {
		return ErrKeyTooLong
	}

	return checkName(string(working))
}

// CheckValue returns if the value is valid.
func CheckValue(value string) error {
	if len(value) > MaxValueLen {
		return ErrValueTooLong
	}
	return checkName(value)
}

func checkName(value string) (err error) {
	valueLen := len(value)
	var state int
	var ch rune
	var width int
	for pos := 0; pos < valueLen; pos += width {
		ch, width = utf8.DecodeRuneInString(value[pos:])
		switch state {
		case 0: //check prefix/suffix
			if !isAlpha(ch) {
				err = exception.NewFromErr(ErrKeyInvalidCharacter).WithMessagef("for: '%s' at: %d", value, pos)
				return
			}
			state = 1
			continue
		case 1:
			if !(isNameSymbol(ch) || ch == BackSlash || isAlpha(ch)) {
				err = exception.NewFromErr(ErrKeyInvalidCharacter).WithMessagef("for: '%s' at: %d", value, pos)
				return
			}
			if pos == valueLen-2 {
				state = 0
			}
			continue
		}
	}
	return
}

func checkDNS(value string) (err error) {
	valueLen := len(value)
	if valueLen == 0 {
		err = exception.Wrap(ErrKeyDNSPrefixEmpty)
		return
	}
	if valueLen > MaxDNSPrefixLen {
		err = exception.Wrap(ErrKeyDNSPrefixTooLong)
		return
	}
	var state int
	var ch rune
	var width int
	for pos := 0; pos < valueLen; pos += width {
		ch, width = utf8.DecodeRuneInString(value[pos:])
		switch state {
		case 0: //check prefix | suffix
			if !isAlpha(ch) {
				return exception.NewFromErr(ErrKeyInvalidCharacter).WithMessagef("for: '%s' at: %d", value, pos)
			}
			state = 1
			continue
		case 1:
			if isNameSymbol(ch) {
				state = 2
				continue
			}
			if !isAlpha(ch) {
				err = exception.NewFromErr(ErrKeyInvalidCharacter).WithMessagef("for: '%s' at: %d", value, pos)
				return
			}
			if pos == valueLen-2 {
				state = 0
			}
			continue
		case 2: // we've hit a dot, dash, or underscore that can't repeat
			if !isAlpha(ch) {
				err = exception.NewFromErr(ErrKeyInvalidCharacter).WithMessagef("for: '%s' at: %d", value, pos)
				return
			}
			if pos == valueLen-2 {
				state = 0
			}

			state = 1
		}
	}
	return nil
}

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isSelectorSymbol(ch rune) bool {
	switch ch {
	case Equal, Bang, OpenParens, CloseParens, Comma:
		return true
	}
	return false
}

func isNameSymbol(ch rune) bool {
	switch ch {
	case Dot, Dash, Underscore:
		return true
	}
	return false
}

func isSymbol(ch rune) bool {
	return (int(ch) >= int(Bang) && int(ch) <= int(ForwardSlash)) ||
		(int(ch) >= int(Colon) && int(ch) <= int(At)) ||
		(int(ch) >= int(OpenBracket) && int(ch) <= int(BackTick)) ||
		(int(ch) >= int(OpenCurly) && int(ch) <= int(Tilde))
}

func isAlpha(ch rune) bool {
	return !isWhitespace(ch) && !unicode.IsControl(ch) && !isSymbol(ch)
}
