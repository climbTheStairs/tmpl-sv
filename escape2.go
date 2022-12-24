//go:build exclude

package main

import (
	"fmt"
	"strings"
)

var escapes = []string{
	`\n`, "\n",
	`\r`, "\r",
	`\t`, "\t",
}
var escaper *strings.Replacer

func init() {
	escaper = strings.NewReplacer(escapes...)
}

func Escape(s string) (string, error) {
	ss := strings.Split(s, `\\`)
	for i, v := range ss {
		v = escaper.Replace(v)
		if len(v) > 0 && v[len(v)-1] == '\\' {
			return "", fmt.Errorf(`unescaped backslash ("\")`)
		}
		if n := strings.Index(v, `\`); n != -1 {
			return "", fmt.Errorf(`invalid escape: \%c`, v[n+1])
		}
		ss[i] = v
	}
	return strings.Join(ss, `\`), nil
}
