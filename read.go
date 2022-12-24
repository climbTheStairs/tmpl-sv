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

	for rowNum := 1; scanner.Scan(); rowNum++ {
		if err := appendRow(&t, scanner.Text()); err != nil {
			return t, fmt.Errorf("row %d: %v", rowNum, err)
		}
	}
	return t, nil
}

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
	for rowNum, row := range t.Body {
		rowNum += 1
		for i, k := range t.Head {
			if row[k], err = Escape(row[k]); err != nil {
				return fmt.Errorf(`row %d: `+
					`column %d "%s": %v`,
					rowNum, i, k, err)
			}
		}
	}
	return nil
}
