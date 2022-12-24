package main

import (
	"html/template"
	"io"
	"strings"
)

// ApplyTmpl reads the template from r,
// executes it on t, and writes the output to w.
func ApplyTmpl(r io.Reader, t *Table, w io.Writer) error {
	fnmap := template.FuncMap{
		"toStrSlice": toStrSlice,
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	tmpl, err := template.New("").Funcs(fnmap).Parse(string(b))
	if err != nil {
		return err
	}
	return tmpl.Execute(w, t.Body)
}

// toStrSlice is meant to be used inside templates.
// It converts a string to a string slice
// by splitting on commas (",").
func toStrSlice(s string) []string {
	return strings.Split(s, ",")
}
