package main

import (
	"bytes"
	"encoding/json"
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

var escaper *strings.Replacer

func init() {
	oldnew := make([]string, 0, len(Escapes)*2)
	for escaped, unescaped := range Escapes {
		oldnew = append(oldnew,
			string(unescaped), "\\"+string(escaped))
	}
	escaper = strings.NewReplacer(oldnew...)
}

// Escape returns a copy of string s
// with each special character disallowed in TSV
// replaced by its two-character escape sequnce.
// See Escapes for specific characters that are escaped
// and their escape sequences.
func Escape(s string) string {
	return escaper.Replace(s)
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

// escapeJson wraps goodJsonMarshal, panicking on errors.
func escapeJson(s string) string {
	b, err := goodJsonMarshal(s)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// goodJsonMarshal is like json.Marshal but good.
// goodJsonMarshal is identical to json.Marshal
// but without its annoying and unasked-for escaping of characters
// that unnecessarily attempts to make the output HTML-safe.
func goodJsonMarshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	// Remove extra newline stupidly added by json.Encoder.Encode.
	b := buf.Bytes()[:buf.Len()-1]
	return b, err
}
