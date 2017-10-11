package selector

import (
	"fmt"
)

const (
	// Dash is a common rune.
	Dash = byte('-')
	// Underscore  is a common rune.
	Underscore = byte('_')
	// Dot is a common rune.
	Dot = byte('.')
	// ForwardSlash is a common rune.
	ForwardSlash = byte('/')
	// BackSlash is a common rune.
	BackSlash = byte('\\')
	// Bang is a common rune.
	Bang = byte('!')
	// Comma is a common rune.
	Comma = byte(',')
	// OpenParens is a common rune.
	OpenParens = byte('(')
	// CloseParens is a common rune.
	CloseParens = byte(')')
	// Equal is a common rune.
	Equal = byte('=')

	// Space is a common rune.
	Space = byte(' ')
	// Tab is a common rune.
	Tab = byte('\t')
	// CarriageReturn is a common rune.
	CarriageReturn = byte('\r')
	// NewLine is a common rune.
	NewLine = byte('\n')
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

	var working []byte
	var state int
	var ch byte

	for pos := 0; pos < keyLen; pos++ {
		ch = key[pos]
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
	var ch byte
	for pos := 0; pos < valueLen; pos++ {
		ch = value[pos]
		switch state {
		case 0: //check prefix | suffix
			if !(isLetter(ch) || isDigit(ch)) {
				err = ErrKeyInvalidCharacter
				return
			}
			state = 1
			continue
		case 1:
			if !(ch == Dot || ch == Dash || ch == Underscore || ch == BackSlash || isLetter(ch) || isDigit(ch)) {
				err = ErrKeyInvalidCharacter
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
		err = ErrKeyDNSPrefixEmpty
		return
	}
	if valueLen > MaxDNSPrefixLen {
		err = ErrKeyDNSPrefixTooLong
		return
	}
	var state int
	var ch byte
	for pos := 0; pos < valueLen; pos++ {
		ch = value[pos]
		switch state {
		case 0: //check prefix | suffix
			if !(isLetter(ch) || isDigit(ch)) {
				return ErrKeyInvalidCharacter
			}
			state = 1
			continue
		case 1:
			if ch == Dot {
				state = 2
				continue
			}
			if !(ch == Dash || ch == Underscore || isLetter(ch) || isDigit(ch)) {
				err = ErrKeyInvalidCharacter
				return
			}
			if pos == valueLen-2 {
				state = 0
			}
			continue
		case 2: // we've hit a dot
			if !(isLetter(ch) || isDigit(ch)) {
				err = ErrKeyInvalidCharacter
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

func isLetter(ch byte) bool {
	return (int(ch) >= int('A') && int(ch) <= int('Z')) || (int(ch) >= int('a') && int(ch) <= int('z'))
}

func isDigit(ch byte) bool {
	return int(ch) >= int('0') && int(ch) <= int('9')
}
