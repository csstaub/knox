package glob

import (
	"testing"
)

func TestGlobPattern(t *testing.T) {
	testCases := []struct {
		pat   string
		cmp   string
		valid bool
	}{
		// Basic patterns
		{"foo", "foo", true},
		{"foo*", "foo", true},
		{"foo*", "foobar", true},
		{"foo*", "foo/bar", false},
		{"foo*/bar", "foo/bar", true},
		{"foo**", "foobar", true},
		{"foo**", "foo/bar", true},
		{"*/bar", "foo/bar", true},
		{"**/bar", "foo/bar", true},
		// Check with . and / in SPIFFE IDs
		{"spiffe://*.example.com/foo*", "spiffe://test.example.com/foo-bar", true},
		{"spiffe://*.example.com/foo*", "spiffe://test.example.com/foo/bar", false},
		{"spiffe://*.example.com/foo**", "spiffe://test.example.com/foo-bar/baz", true},
		{"spiffe://*.example.com/foo**", "spiffe://test.example.com/foo/bar/baz", true},
		{"spiffe://*.example.com/foo**", "spiffe://example.com/foo/bar/baz", false},
		// Check that regex chars escaped
		{"?", "?", true},
		{"?", "!", false},
		{"\\", "\\", true},
		{"\\w", "\\w", true},
		{"\\w", "x", false},
		{".*", ".*", true},
		{".*", "x", false},
		{"(x|y)", "(x|y)", true},
		{"(x|y)", "x", false},
		{"(x|y)", "y", false},
	}

	for _, tc := range testCases {
		pattern, err := NewPattern(tc.pat)
		if err != nil {
			t.Errorf("unexpected invalid pattern: %s; error: %v", tc.pat, err)
			continue
		}
		matches := pattern.MatchString(tc.cmp)
		if matches != tc.valid {
			if matches {
				t.Errorf("pattern %s matches %s, but should not", tc.pat, tc.cmp)
			} else {
				t.Errorf("pattern %s doesn't match %s, but should", tc.pat, tc.cmp)
			}
		}
	}
}

func TestGlobPatternInvalid(t *testing.T) {
	invalidPatterns := []string{
		"",
		"***",
		"x***",
		"***z",
		"x***z",
	}

	for _, pat := range invalidPatterns {
		_, err := NewPattern(pat)
		if err == nil {
			t.Errorf("pattern %s was supposed to be invalid, but yielded no error", pat)
		}
	}
}
