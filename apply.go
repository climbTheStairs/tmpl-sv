package main

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// ApplyTmpl reads the template from the file called fname
// and executes it on t, writing the output to out.
func ApplyTmpl(fname string, t *Table, out io.Writer) error {
	fnmap := template.FuncMap{
		"toStrSlice": toStrSlice,
	}
	tmpl, err := template.
		New(filepath.Base(fname)).
		Funcs(fnmap).
		ParseFiles(fname)
	if err != nil {
		return err
	}
	return tmpl.Execute(out, t.Body)
}

// toStrSlice is meant to be used inside templates.
// It converts a string to a string slice
// by splitting on commas (",").
func toStrSlice(s string) []string {
	return strings.Split(s, ",")
}
