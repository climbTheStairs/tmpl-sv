package main

import (
	"html/template"
	"io"
	"strings"
)

var funcMap = template.FuncMap{
	"toStrSlice": toStrSlice,
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

// toStrSlice is meant to be used inside templates.
// It converts a string to a string slice
// by splitting on commas (",").
func toStrSlice(s string) []string {
	return strings.Split(s, ",")
}
