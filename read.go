package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func readTsv(f io.Reader) ([]map[string]string, []string, error) {
	scanner := bufio.NewScanner(f)
	scanner.Split(scanPosixLines)

	if ok := scanner.Scan(); !ok {
		return nil, nil, fmt.Errorf("cannot read initial line")
	}
	head := strings.Split(scanner.Text(), "\t")

	t := make([]map[string]string, 0)
	for lnum := 2; scanner.Scan(); lnum++ {
		row, err := createRow(scanner.Text(), head)
		if err != nil {
			return t, head, fmt.Errorf("line %d: %v", lnum, err)
		}
		t = append(t, row)
	}
	return t, head, nil
}

func createRow(s string, head []string) (map[string]string, error) {
	cols := strings.Split(s, "\t")
	if len(cols) != len(head) {
		return nil, fmt.Errorf("invalid number of columns")
	}
	row := make(map[string]string)
	for i, v := range head {
		row[v] = cols[i]
	}
	return row, nil
}

func escapeTable(t []map[string]string, head []string) ([]map[string]string, error) {
	for i := range t {
		var err error
		lnum := i + 2
		if t[i], err = escapeRow(t[i], head); err != nil {
			return t, fmt.Errorf("line %d: %v", lnum, err)
		}
	}
	return t, nil
}

func escapeRow(row map[string]string, head []string) (map[string]string, error) {
	for i, k := range head {
		var err error
		if row[k], err = escape(row[k]); err != nil {
			return row, fmt.Errorf(`column %d "%s": %v`, i, k, err)
		}
	}
	return row, nil
}
