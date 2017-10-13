package selector

import (
	"strings"
	"unicode/utf8"
)

const (
	// OpEquals is an operator.
	OpEquals = "="
	// OpDoubleEquals is an operator.
	OpDoubleEquals = "=="
	// OpNotEquals is an operator.
	OpNotEquals = "!="
	// OpIn is an operator.
	OpIn = "in"
	// OpNotIn is an operator.
	OpNotIn = "notin"
)

// Lexer is the working engine of the semantic extraction for a selector.
// It lets us work through a string with a cursor, with an optional mark we can refer back  to.
type Lexer struct {
	// s stores the string to be tokenized
	s string
	// pos is the position currently tokenized
	pos int
	// m is an optional mark
	m int
}

// Lex does the actual parsing.
func (l *Lexer) Lex() (Selector, error) {
	l.s = strings.TrimSpace(l.s)
	if len(l.s) == 0 {
		return nil, ErrEmptySelector
	}

	var b rune
	var selector Selector
	var err error
	var op string

	// loop over "clauses"
	for {

		// sniff the !haskey form
		b = l.current()
		if b == Bang {
			l.advance() // we aren't going to use the '!'
			selector = l.lift(selector, l.notHasKey(l.readWord()))
			if l.done() {
				break
			}
			continue
		}

		// we're done peeking the first char
		key := l.readWord()

		l.mark()
		b = l.skipToComma()
		if b == Comma || l.isTerminator(b) || l.done() {
			selector = l.lift(selector, l.hasKey(key))
			l.advance()
			if l.done() {
				break
			}
			continue
		} else {
			l.popMark()
		}

		op, err = l.readOp()
		if err != nil {
			return nil, err
		}

		var subSelector Selector
		switch op {
		case OpEquals, OpDoubleEquals:
			subSelector, err = l.equals(key)
			if err != nil {
				return nil, err
			}
			selector = l.lift(selector, subSelector)
		case OpNotEquals:
			subSelector, err = l.notEquals(key)
			if err != nil {
				return nil, err
			}
			selector = l.lift(selector, subSelector)
		case OpIn:
			subSelector, err = l.in(key)
			if err != nil {
				return nil, err
			}
			selector = l.lift(selector, subSelector)
		case OpNotIn:
			subSelector, err = l.notIn(key)
			if err != nil {
				return nil, err
			}
			selector = l.lift(selector, subSelector)
		default:
			return nil, ErrInvalidOperator
		}

		b = l.skipToComma()
		if b == Comma {
			l.advance()
			if l.done() {
				break
			}
			continue
		}

		// these two are effectively the same
		if l.isTerminator(b) || l.done() {
			break
		}

		return nil, ErrInvalidSelector
	}

	err = selector.Validate()
	if err != nil {
		return nil, err
	}

	return selector, nil
}

// lift starts grouping selectors into a high level `and`, returning the aggregate selector.
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

func (l *Lexer) equals(key string) (Selector, error) {
	value := l.readWord()
	return Equals{Key: key, Value: value}, nil
}

func (l *Lexer) notEquals(key string) (Selector, error) {
	value := l.readWord()
	return NotEquals{Key: key, Value: value}, nil
}

func (l *Lexer) in(key string) (Selector, error) {
	csv, err := l.readCSV()
	if err != nil {
		return nil, err
	}
	return In{Key: key, Values: csv}, nil
}

func (l *Lexer) notIn(key string) (Selector, error) {
	csv, err := l.readCSV()
	if err != nil {
		return nil, err
	}
	return NotIn{Key: key, Values: csv}, nil
}

// done indicates the cursor is past the usable length of the string.
func (l *Lexer) done() bool {
	return l.pos == len(l.s)
}

// mark sets a mark at the current position.
func (l *Lexer) mark() {
	l.m = l.pos
}

// popMark moves the cursor back to the previous mark.
func (l *Lexer) popMark() {
	if l.m > 0 {
		l.pos = l.m
	}
	l.m = 0
}

// read returns the rune currently lexed, and advances the position.
func (l *Lexer) read() (r rune) {
	var width int
	if l.pos < len(l.s) {
		r, width = utf8.DecodeRuneInString(l.s[l.pos:])
		l.pos += width
	}
	return r
}

// current returns the rune at the current position.
func (l *Lexer) current() (r rune) {
	r, _ = utf8.DecodeRuneInString(l.s[l.pos:])
	return
}

// advance moves the cursor forward one rune.
func (l *Lexer) advance() {
	if l.pos < len(l.s) {
		_, width := utf8.DecodeRuneInString(l.s[l.pos:])
		l.pos += width
	}
}

// unread moves the cursor back a rune.
func (l *Lexer) prev() {
	if l.pos > 0 {
		l.pos--
	}
}

// readOp reads a valid operator.
// valid operators include:
// [ =, ==, !=, in, notin ]
// errors if it doesn't read one of the above, or there is another structural issue.
func (l *Lexer) readOp() (string, error) {
	// skip preceding whitespace
	l.skipWhiteSpace()

	var state int
	var ch rune
	var op []rune
	for {
		ch = l.current()

		switch state {
		case 0: // initial state, determine what op we're reading for
			if ch == Equal {
				state = 1
				break
			}
			if ch == Bang {
				state = 2
				break
			}
			if ch == 'i' {
				state = 6
				break
			}
			if ch == 'n' {
				state = 7
				break
			}
			return "", ErrInvalidOperator
		case 1: // =
			if l.isWhitespace(ch) || l.isAlpha(ch) || ch == Comma {
				return string(op), nil
			}
			if ch == Equal {
				op = append(op, ch)
				l.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 2: // !
			if ch == Equal {
				op = append(op, ch)
				l.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 6: // in
			if ch == 'n' {
				op = append(op, ch)
				l.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 7: // o
			if ch == 'o' {
				state = 8
				break
			}
			return "", ErrInvalidOperator
		case 8: // t
			if ch == 't' {
				state = 9
				break
			}
			return "", ErrInvalidOperator
		case 9: // i
			if ch == 'i' {
				state = 10
				break
			}
			return "", ErrInvalidOperator
		case 10: // n
			if ch == 'n' {
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
	// skip preceding whitespace
	l.skipWhiteSpace()

	var word []rune
	var ch rune
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

func (l *Lexer) readCSV() (results []string, err error) {
	// skip preceding whitespace
	l.skipWhiteSpace()

	var word []rune
	var ch rune
	var state int

	for {
		ch = l.current()

		if l.done() {
			err = ErrInvalidSelector
			return
		}

		switch state {

		case 0: // leading paren

			if ch == OpenParens {
				state = 2 // spaces or alphas
				l.advance()
				continue
			}

		case 1: // alphas (in word)

			if ch == Comma {
				if len(word) > 0 {
					results = append(results, string(word))
					word = nil
				}
				state = 2 // from comma
				l.advance()
				continue
			}

			if ch == CloseParens {
				if len(word) > 0 {
					results = append(results, string(word))
				}
				l.advance()
				return
			}

			if l.isWhitespace(ch) {
				state = 3
				l.advance()
				continue
			}

			if !l.isAlpha(ch) {
				err = ErrInvalidSelector
				return
			}

			word = append(word, ch)
			l.advance()
			continue

		case 2: //whitespace after symbol

			if ch == CloseParens {
				l.advance()
				return
			}

			if l.isWhitespace(ch) {
				l.advance()
				continue
			}

			if ch == Comma {
				l.advance()
				continue
			}

			if l.isAlpha(ch) {
				state = 1
				continue
			}

			err = ErrInvalidSelector
			return

		case 3: //whitespace after alpha

			if ch == CloseParens {
				if len(word) > 0 {
					results = append(results, string(word))
				}
				l.advance()
				return
			}

			if l.isWhitespace(ch) {
				l.advance()
				continue
			}

			if ch == Comma {
				if len(word) > 0 {
					results = append(results, string(word))
					word = nil
				}
				l.advance()
				state = 2
				continue
			}

			err = ErrInvalidSelector
			return

		}
	}
}

func (l *Lexer) skipWhiteSpace() {
	if l.done() {
		return
	}
	var ch rune
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

func (l *Lexer) skipToComma() (ch rune) {
	if l.done() {
		return
	}
	for {
		ch = l.current()
		if ch == Comma {
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
func (l *Lexer) isWhitespace(ch rune) bool {
	return ch == Space || ch == Tab || ch == CarriageReturn || ch == NewLine
}

// isSpecialSymbol returns if the ch can be a token.
func (l *Lexer) isSpecialSymbol(ch rune) bool {
	return isSelectorSymbol(ch)
}

// isTerminator returns if we've reached the end of the string
func (l *Lexer) isTerminator(ch rune) bool {
	return ch == 0
}

func (l *Lexer) isAlpha(ch rune) bool {
	return isAlpha(ch)
}
