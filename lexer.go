package selector

import (
	"fmt"
	"strings"
)

var (
	// ErrInvalidOperator is returned if the operator is invalid.
	ErrInvalidOperator = fmt.Errorf("invalid operator")

	// ErrInvalidSelector is returned if there is a structural issue with the selector.
	ErrInvalidSelector = fmt.Errorf("invalid selector")
)

// Lexer represents the Lexer struct for label selector.
// It contains necessary informationt to tokenize the input string
type Lexer struct {
	// s stores the string to be tokenized
	s string
	// pos is the position currently tokenized
	pos int
}

// Lex returns a pair of Token and the literal
// literal is meaningfull only for IdentifierToken token
func (l *Lexer) Lex() (Selector, error) {
	l.s = strings.TrimSpace(l.s)
	if len(l.s) == 0 {
		return nil, fmt.Errorf("query is empty")
	}

	var b byte
	var selector Selector
	var err error
	var op string

	// loop over "clauses"
	for {

		// sniff the !haskey form
		b = l.current()
		if b == '!' {
			l.advance() // we aren't going to use the '!'
			selector = l.lift(selector, l.notHasKey(l.readWord()))
			break
		}

		// we're done peeking the first char
		key := l.readWord()

		if l.done() {
			selector = l.lift(selector, l.hasKey(key))
			break
		}

		op, err = l.readOp()
		if err != nil {
			return nil, err
		}

		switch op {
		case "=", "==":
			selector = l.lift(selector, l.equals(key))
		case "!=":
			selector = l.lift(selector, l.notEquals(key))
		case "in":
			selector = l.lift(selector, l.in(key))
		case "notin":
			selector = l.lift(selector, l.notIn(key))
		default:
			return nil, ErrInvalidOperator
		}

		b = l.skipToComma()
		if b == byte(',') {
			l.advance()
			if l.done() {
				break
			}
			continue
		}

		// this is the same thing.
		if l.isTerminator(b) || l.done() {
			break
		}

		println("bad selector; not finished, not comma")
		return nil, ErrInvalidSelector
	}

	return selector, nil
}

func (l *Lexer) lift(current, next Selector) Selector {
	if current == nil {
		return next
	}
	if typed, isTyped := current.(And); isTyped {
		return append(typed, next)
	}
	return And([]Selector{current, next})
}

func (l *Lexer) hasKey(key string) Selector {
	return HasKey(key)
}

func (l *Lexer) notHasKey(key string) Selector {
	return NotHasKey(key)
}

func (l *Lexer) equals(key string) Selector {
	value := l.readWord()
	return Equals{Key: key, Value: value}
}

func (l *Lexer) notEquals(key string) Selector {
	value := l.readWord()
	return NotEquals{Key: key, Value: value}
}

func (l *Lexer) in(key string) Selector {
	return In{Key: key, Values: l.readCSV()}
}

func (l *Lexer) notIn(key string) Selector {
	return NotIn{Key: key, Values: l.readCSV()}
}

// done indicates the cursor is past the usable length of the string.
func (l *Lexer) done() bool {
	return l.pos == len(l.s)
}

// read return the character currently lexed
// increment the position and check the buffer overflow
func (l *Lexer) read() (b byte) {
	if l.pos < len(l.s) {
		b = l.s[l.pos]
		l.pos++
	}
	return b
}

// current returns the byte a the current position.
func (l *Lexer) current() byte {
	return l.s[l.pos]
}

// advance moves the cursor forward one character.
func (l *Lexer) advance() {
	if l.pos < len(l.s) {
		l.pos++
	}
}

// unread moves the cursor back a character.
func (l *Lexer) prev() {
	if l.pos > 0 {
		l.pos--
	}
}

// readOp reads a valid operator.
// valid operators include:
// [ =, ==, !=, in, notin ]
// errors if it doesn't read one of the above.
func (l *Lexer) readOp() (string, error) {
	// skip preceeding whitespace
	l.skipWhiteSpace()

	var state int
	var ch byte
	var op []byte
	for {
		ch = l.current()

		switch state {
		case 0: // initial state, determine what op we're reading for
			if ch == byte('=') {
				state = 1
				break
			}
			if ch == byte('!') {
				state = 2
				break
			}
			if ch == byte('i') {
				state = 6
				break
			}
			if ch == byte('n') {
				state = 7
				break
			}
			return "", ErrInvalidOperator
		case 1: // =
			if l.isWhitespace(ch) || l.isAlpha(ch) {
				return string(op), nil
			}
			if ch == byte('=') {
				op = append(op, ch)
				l.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 2: // !
			if ch == byte('=') {
				op = append(op, ch)
				l.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 6: // in
			if ch == byte('n') {
				op = append(op, ch)
				l.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 7: // o
			if ch == byte('o') {
				state = 8
				break
			}
			return "", ErrInvalidOperator
		case 8: // t
			if ch == byte('t') {
				state = 9
				break
			}
			return "", ErrInvalidOperator
		case 9: // i
			if ch == byte('i') {
				state = 10
				break
			}
			return "", ErrInvalidOperator
		case 10: // n
			if ch == byte('n') {
				op = append(op, ch)
				l.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		}

		op = append(op, ch)
		l.advance()

		if l.done() {
			return string(op), nil
		}
	}
}

// readWord skips whitespace, then reads a word until whitespace or a token.
// it will leave the cursor on the next char after the word, i.e. the space or token.
func (l *Lexer) readWord() string {
	// skip preceeding whitespace
	l.skipWhiteSpace()

	var word []byte
	var ch byte
	for {
		ch = l.current()

		if l.isWhitespace(ch) {
			return string(word)
		}
		if l.isSpecialSymbol(ch) {
			return string(word)
		}
		word = append(word, ch)
		l.advance()

		if l.done() {
			return string(word)
		}
	}
}

func (l *Lexer) readCSV() (results []string) {
	// skip preceeding whitespace
	l.skipWhiteSpace()

	var word []byte
	var ch byte
	for {
		ch = l.current()
		if ch == byte(')') {
			if len(word) > 0 {
				results = append(results, string(word))
			}
			l.advance()
			return
		}

		if ch == byte('(') || l.isWhitespace(ch) {
			l.advance()
			continue
		}

		if ch == byte(',') {
			results = append(results, string(word))
			word = []byte{}
			l.advance()
			continue
		}

		word = append(word, ch)
		l.advance()
		if l.done() {
			if len(word) > 0 {
				results = append(results, string(word))
			}
			return
		}
	}
}

func (l *Lexer) skipWhiteSpace() {
	if l.done() {
		return
	}
	var ch byte
	for {
		ch = l.current()
		if !l.isWhitespace(ch) {
			return
		}
		l.advance()
		if l.done() {
			return
		}
	}
}

func (l *Lexer) skipToComma() (ch byte) {
	if l.done() {
		return
	}
	for {
		ch = l.current()
		if ch == byte(',') {
			return
		}
		if !l.isWhitespace(ch) {
			return
		}
		l.advance()
		if l.done() {
			return
		}
	}
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func (l *Lexer) isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

// isSpecialSymbol detect if the character ch can be an operator
func (l *Lexer) isSpecialSymbol(ch byte) bool {
	switch ch {
	case '=', '!', '(', ')', ',', '>', '<':
		return true
	}
	return false
}

// isTerminator returns if we've reached the end of the string
func (l *Lexer) isTerminator(ch byte) bool {
	return ch == 0
}

// this needs a test la.
func (l *Lexer) isAlpha(ch byte) bool {
	return int(ch) >= int('A') && int(ch) <= int('z')
}
