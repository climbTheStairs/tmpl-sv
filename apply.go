package main

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// TODO: Would passing Table instead of *Table
// be significantly more expensive?
func applyTmpl(fname string, t *Table, out io.Writer) error {
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

func toStrSlice(s string) []string {
	return strings.Split(s, ",")
}
