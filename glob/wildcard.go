package glob

import (
	"bytes"
	"errors"
	"regexp"
)

var (
	// Specifies characters to exclude when matching with '*' in a pattern.
	separators = []rune{'/', '.'}
)

// New creates a new pattern instance.
func NewPattern(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, errors.New("invalid pattern")
	}

	// Add ^ to make sure we match full string every time
	var regex bytes.Buffer
	regex.WriteRune('^')

	asterisksSeen := 0
	for _, r := range pattern {
		switch r {
		case '*':
			asterisksSeen += 1
		default:
			err := handleAsterisks(asterisksSeen, &regex)
			if err != nil {
				return nil, err
			}
			asterisksSeen = 0

			// Match exact rune (escaped)
			regex.WriteString(regexp.QuoteMeta(string(r)))
		}
	}

	// Handle any trailing asterisks
	err := handleAsterisks(asterisksSeen, &regex)
	if err != nil {
		return nil, err
	}

	// Don't forget final terminator to match full string
	regex.WriteRune('$')

	compiled, err := regexp.Compile(regex.String())
	if err != nil {
		return nil, err
	}

	return compiled, nil
}

// Generate the appropriate regex pattern for a given number of asterisks.
func handleAsterisks(asterisksSeen int, regex *bytes.Buffer) error {
	if asterisksSeen == 1 {
		// Single asterisk -- don't match separators.
		regex.WriteString("[^")
		for _, sep := range separators {
			regex.WriteRune(sep)
		}
		regex.WriteString("]*")
	} else if asterisksSeen == 2 {
		// Double asterisk -- match all including separators.
		regex.WriteString(".*")
	} else if asterisksSeen > 2 {
		// Triple or more asterisks -- invalid.
		return errors.New("invalid glob pattern: saw more than 2 asterisks in a row")
	}
	return nil
}
