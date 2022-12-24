package main

import (
	"fmt"
	"strings"
)

var escapes = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\\': '\\',
}

// Escape returns a copy of string s
// with each non-overlapping two-character substring beginning with '\\'
// replaced by the character corresponding to that escape sequence.
// If s contains invalid escapes or unescaped backslashes ("\"),
// Escape makes no further replacements and returns a non-nil error.
func Escape(s string) (string, error) {
	var b strings.Builder
	b.Grow(len(s))
	start := 0
	for i := 0; i < len(s); {
		if s[i] != '\\' {
			// s[i] does not start an escape;
			// keep moving forward in s until it does.
			i += 1
			continue
		}
		if i+1 == len(s) {
			// s[i] is the last character in s
			// and therefore cannot start an escape.
			return b.String(),
				fmt.Errorf(`unescaped backslash ("\")`)
		}
		escaped, ok := escapes[s[i+1]]
		if !ok {
			return b.String(),
				fmt.Errorf(`invalid escape: \%c`, s[i+1])
		}
		b.WriteString(s[start:i])
		b.WriteByte(escaped)
		// Move forward in s by 2,
		// to pass the two-character escape substring.
		i += 2
		start = i
	}
	// The last substring of s that does not end with an escape
	// must be written if it exists.
	b.WriteString(s[start:])
	return b.String(), nil
}
