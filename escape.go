package main

import (
	"fmt"
	"strings"
)

// Escapes is a map of the second character
// of each two-character escape sequence
// (i.e. the character that isn't a backslash ('\\'))
// to the special character
// represented by the whole escape sequence.
var Escapes = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\\': '\\',
}

// Unescape returns a copy of string s
// with each non-overlapping two-character escape sequence
// beginning with a backslash ('\\')
// replaced by the special character represented by the escape.
// See Escapes for specific escape sequences
// and their corresponding characters.
// If s contains invalid escapes or unescaped backslashes,
// Unescape makes no further replacements and returns a non-nil error.
func Unescape(s string) (string, error) {
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
		unescaped, ok := Escapes[s[i+1]]
		if !ok {
			return b.String(),
				fmt.Errorf(`invalid escape: \%c`, s[i+1])
		}
		b.WriteString(s[start:i])
		b.WriteByte(unescaped)
		// Move forward in s by 2,
		// to pass by the two-character escape sequence.
		i += 2
		start = i
	}
	// The last substring of s that does not end with an escape
	// must be written if it exists.
	b.WriteString(s[start:])
	return b.String(), nil
}
