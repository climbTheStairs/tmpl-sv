package main

import (
	"html/template"
	"io"
	"strconv"
	"strings"
)

var funcMap = template.FuncMap{
	"int":     strconv.Atoi,
	"split":   strings.Split,
	"add":     func(a, b int) int { return a + b },
	"mkslice": func(a ...string) []string { return a },
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
