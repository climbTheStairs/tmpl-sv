package main

import (
	"io"
	"strconv"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"int":   strconv.Atoi,
	"split": strings.Split,
	"str":   strconv.Itoa,
	"add":   func(a, b int) int {
		return a + b
	},
	"map":   func(props ...string) map[string]string {
		if len(props) % 2 != 0 {
			return nil
		}
		m := make(map[string]string)
		for i := 0; i < len(props); i += 2 {
			m[props[i]] = props[i+1]
		}
		return m
	},
}

// ReadTemplate reads from r
// and returns a template
// with a function map added.
func ReadTemplate(r io.Reader) (*template.Template, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return template.New("").Funcs(funcMap).Parse(string(b))
}
