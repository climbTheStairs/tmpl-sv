package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ReadTsv creates and returns a new table
// using TSV data read from f.
func ReadTsv(f io.Reader, esc bool) (*Table, error) {
	var t *Table
	var err error

	scanner := bufio.NewScanner(f)
	scanner.Split(scanPosixLines)

	if ok := scanner.Scan(); !ok {
		return t, fmt.Errorf("cannot read initial line")
	}
	t.Head = strings.Split(scanner.Text(), "\t")

	for i := 1; scanner.Scan(); i++ {
		fields := strings.Split(scanner.Text(), "\t")
		if fields, err = unescapeFields(t, fields, esc); err != nil {
			return t, fmt.Errorf("row %d: %v", i, err)
		}
		if err := t.AppendRow(fields); err != nil {
			return t, fmt.Errorf("row %d: %v", i, err)
		}
	}
	return t, nil
}

// unescapeFields is a helper function for ReadTsv.
// If esc is true, unescapeFields replaces each field in fields
// with the output of calling Unescape on that field.
// If any field causes Unescape to return an error,
// unescapeFields returns that wrapped error.
// If esc is false, unescapeFields returns fields unchanged.
func unescapeFields(t *Table, fields []string, esc bool) ([]string, error) {
	if !esc {
		return fields, nil
	}
	for i, v := range fields {
		unescaped, err := Unescape(v)
		if err != nil {
			return fields, fmt.Errorf(`column %d "%s": %v`,
				i+1, t.Head[i], err)
		}
		fields[i] = unescaped
	}
	return fields, nil
}

// AppendRow creates a row from fields
// and appends it to table t.
func (t *Table) AppendRow(fields []string) error {
	if len(fields) != len(t.Head) {
		return fmt.Errorf("invalid number of columns")
	}
	row := make(map[string]string)
	for i, k := range t.Head {
		row[k] = fields[i]
	}
	t.Body = append(t.Body, row)
	return nil
}
