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

func escape(s string) (string, error) {
	var b strings.Builder
	b.Grow(len(s))

	start := 0
	for i := 0; i < len(s)-1; {
		if s[i] != '\\' {
			i += 1
			continue
		}
		escaped, ok := escapes[s[i+1]]
		if !ok {
			return b.String(),
				fmt.Errorf(`invalid escape: \%c`, s[i+1])
		}
		b.WriteString(s[start:i])
		b.WriteByte(escaped)
		i += 2
		start = i
	}
	b.WriteString(s[start:])

	if len(s) > 0 && s[len(s)-1] == '\\' && start != len(s) {
		return b.String(), fmt.Errorf(`unescaped backslash ("\")`)
	}

	return b.String(), nil
}
