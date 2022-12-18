package main

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

func applyTmpl(fname string, d []map[string]string, out io.Writer) error {
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
	return tmpl.Execute(out, d)
}

func toStrSlice(s string) []string {
	return strings.Split(s, ",")
}
