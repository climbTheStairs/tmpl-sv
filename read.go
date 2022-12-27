package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ReadTsv creates and returns a new table
// using TSV data read from f.
func ReadTsv(f io.Reader) (Table, error) {
	t := Table{}

	scanner := bufio.NewScanner(f)
	scanner.Split(scanPosixLines)

	if ok := scanner.Scan(); !ok {
		return t, fmt.Errorf("cannot read initial line")
	}
	t.Head = strings.Split(scanner.Text(), "\t")

	for i := 1; scanner.Scan(); i++ {
		cols := strings.Split(scanner.Text(), "\t")
		if err := t.AppendRow(cols); err != nil {
			return t, fmt.Errorf("row %d: %v", i, err)
		}
	}
	return t, nil
}

// AppendRow creates a row from cols
// and appends it to table t.
func (t *Table) AppendRow(cols []string) error {
	if len(cols) != len(t.Head) {
		return fmt.Errorf("invalid number of columns")
	}
	row := make(map[string]string)
	for i, v := range t.Head {
		row[v] = cols[i]
	}
	t.Body = append(t.Body, row)
	return nil
}

// UnescapeTable replaces each field in t
// with the output of calling Unescape on that field.
// If any fields contain invalid escapes or unescaped backslashes ('\\').
// UnescapeTable makes no further replacements and returns a non-nil error.
func (t *Table) UnescapeTable() error {
	for i, row := range t.Body {
		i += 1
		for j, k := range t.Head {
			j += 1
			unescaped, err := Unescape(row[k])
			if err != nil {
				return fmt.Errorf(`row %d: column %d "%s": %v`,
					i, j, k, err)
			}
			row[k] = unescaped
		}
	}
	return nil
}
