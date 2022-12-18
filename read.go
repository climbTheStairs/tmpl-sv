package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func readTsv(f io.Reader) (Table, error) {
	t := Table{}

	scanner := bufio.NewScanner(f)
	scanner.Split(scanPosixLines)

	if ok := scanner.Scan(); !ok {
		return t, fmt.Errorf("cannot read initial line")
	}
	t.Head = strings.Split(scanner.Text(), "\t")

	for lnum := 2; scanner.Scan(); lnum++ {
		if err := appendRow(&t, scanner.Text()); err != nil {
			return t, fmt.Errorf("line %d: %v", lnum, err)
		}
	}
	return t, nil
}

// TODO: What would be different if I used instead
// a func (t *Table) method?
func appendRow(t *Table, s string) error {
	cols := strings.Split(s, "\t")
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

func escapeTable(t *Table) error {
	var err error
	for lnum, row := range t.Body {
		lnum += 2
		for i, k := range t.Head {
			if row[k], err = escape(row[k]); err != nil {
				return fmt.Errorf(`line %d: `+
					`column %d "%s": %v`,
					lnum, i, k, err)
			}
		}
	}
	return nil
}
