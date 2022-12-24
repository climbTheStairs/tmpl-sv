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

	for rowNum := 1; scanner.Scan(); rowNum++ {
		cols := strings.Split(scanner.Text(), "\t")
		if err := appendRow(&t, cols); err != nil {
			return t, fmt.Errorf("row %d: %v", rowNum, err)
		}
	}
	return t, nil
}

// AppendRow creates a row from cols
// and appends it to table t.
func AppendRow(t *Table, cols []string) error {
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

// EscapeTable replaces each cell in t
// with the output of calling Escape on that cell.
// If any cells contain invalid escapes or unescaped backslashes ("\").
// EscapeTable makes no further replacements and returns a non-nil error.
func EscapeTable(t *Table) error {
	for rowNum, row := range t.Body {
		rowNum += 1
		for colNum, colName := range t.Head {
			colNum += 1
			escaped, err := Escape(row[colName])
			if err != nil {
				return fmt.Errorf(`row %d: column %d "%s": %v`,
					rowNum, colNum, colName, err)
			}
			row[colName] = escaped
		}
	}
	return nil
}

// escapeIf wraps escape, returning s
// escaped if esc is true, unchanged otherwise.
func escapeIf(s string, esc bool) (string, error) {
	if esc {
		return escape(s)
	}
	return s, nil
}

func escape(s string) (string, error) {
	ss := strings.Split(s, `\\`)
	for i, v := range ss {
		ss[i] = escaper.Replace(v)
		if strings.Contains(ss[i], `\`) {
			return "", fmt.Errorf("...")
		}
	}
	return strings.Join(ss, `\`), nil
}
