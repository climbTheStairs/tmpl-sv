package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ReadTsv creates and returns a new table
// using TSV data read from f.
func ReadTsv(f io.Reader) (*Table, error) {
	var t *Table

	scanner := bufio.NewScanner(f)
	scanner.Split(scanPosixLines)

	if ok := scanner.Scan(); !ok {
		return t, fmt.Errorf("cannot read initial line")
	}
	t.Head = strings.Split(scanner.Text(), "\t")

	for i := 1; scanner.Scan(); i++ {
		fields := strings.Split(scanner.Text(), "\t")
		if err := t.AppendRow(fields); err != nil {
			return t, fmt.Errorf("row %d: %v", i, err)
		}
	}
	return t, nil
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

// UnescapeTable replaces each field in t
// with the output of calling Unescape on that field.
// If any field causes Unescape to return an error,
// UnescapeTable returns that wrapped error.
func (t *Table) UnescapeTable() error {
	for i, row := range t.Body {
		for j, k := range t.Head {
			unescaped, err := Unescape(row[k])
			if err != nil {
				return fmt.Errorf(`row %d: column %d "%s": %v`,
					i+1, j+1, k, err)
			}
			row[k] = unescaped
		}
	}
	return nil
}
