package selector

import (
	"fmt"
	"strings"
)

// NotIn returns if a key does not match a set of values.
type NotIn struct {
	Key    string
	Values []string
}

// Matches returns the selector result.
func (ni NotIn) Matches(labels Labels) bool {
	if value, hasValue := labels[ni.Key]; hasValue {
		for _, iv := range ni.Values {
			if iv == value {
				return false
			}
		}
	}
	return true
}

// Validate validates the selector.
func (ni NotIn) Validate() (err error) {
	err = CheckLabel(ni.Key)
	if err != nil {
		return
	}
	for _, v := range ni.Values {
		err = CheckValue(v)
		if err != nil {
			return
		}
	}
	return
}

// String returns a string representation of the selector.
func (ni NotIn) String() string {
	return fmt.Sprintf("%s notin (%s)", ni.Key, strings.Join(ni.Values, ", "))
}
