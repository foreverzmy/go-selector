package selector

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	// RuneDash is a common rune.
	RuneDash = rune('-')
	// RuneUnderscore  is a common rune.
	RuneUnderscore = rune('_')
	// RuneDot         is a common rune.
	RuneDot = rune('.')
)

// CheckLabel validates a label
func CheckLabel(label string) error {
	// check for '/', if more than one, error
	if strings.Contains(label, "/") {
		parts := strings.Split(label, "/")
		if len(parts) > 2 {
			return fmt.Errorf("label contains more than one slash separator")
		}

		if err := checkLabelPrefix(parts[0]); err != nil {
			return err
		}

		return checkName(parts[1])
	}
	return checkName(label)
}

// CheckValue returns if the value is valid.
func CheckValue(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("value is empty")
	}
	if len(value) > 63 {
		return fmt.Errorf("value is too long")
	}
	return checkName(value)
}

func checkLabelPrefix(prefix string) error {
	if len(prefix) == 0 {
		return fmt.Errorf("label prefix is empty")
	}
	if len(prefix) > 253 {
		return fmt.Errorf("label prefix is too long; must be less than 253 characters")
	}

	return checkName(prefix)
}

func checkLabelSuffix(suffix string) error {
	if len(suffix) == 0 {
		return fmt.Errorf("label suffix is empty")
	}
	if len(suffix) > 63 {
		return fmt.Errorf("label suffix is too long; must be less than 63 characters")
	}

	return checkName(suffix)
}

func checkName(value string) error {
	var state int
	for index, c := range value {
		switch state {
		case 0:
			if !(unicode.IsLetter(c) || unicode.IsDigit(c)) {
				return fmt.Errorf("invalid character (%v) at %d", c, index)
			}
			state = 1
			continue
		case 1:
			if !(c == RuneDash ||
				c == RuneUnderscore ||
				c == RuneDot ||
				unicode.IsLetter(c) ||
				unicode.IsDigit(c)) {

				return fmt.Errorf("invalid character (%v) at %d", c, index)
			}
			if index == len(value)-2 {
				state = 0
			}
			continue
		}
	}
	return nil
}
