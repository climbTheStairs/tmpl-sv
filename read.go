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

// Unescape returns a copy of string s
// with each non-overlapping two-character escape sequence
// beginning with a backslash ('\\')
// replaced by the special character represented by the escape.
// See Escapes for specific escape sequences
// and their corresponding characters.
// If s contains invalid escapes or unescaped backslashes,
// Unescape makes no further replacements and returns a non-nil error.
func Unescape(s string) (string, error) {
	var b strings.Builder
	b.Grow(len(s))
	start := 0
	for i := 0; i < len(s); {
		if s[i] != '\\' {
			// s[i] does not start an escape;
			// keep moving forward in s until it does.
			i += 1
			continue
		}
		if i+1 == len(s) {
			// s[i] is the last character in s
			// and therefore cannot start an escape.
			return b.String(), fmt.Errorf(`unescaped backslash ("\")`)
		}
		unescaped, ok := Escapes[s[i+1]]
		if !ok {
			return b.String(), fmt.Errorf(`invalid escape: \%c`, s[i+1])
		}
		b.WriteString(s[start:i])
		b.WriteByte(unescaped)
		// Move forward in s by 2,
		// to pass by the two-character escape sequence.
		i += 2
		start = i
	}
	// The last substring of s that does not end with an escape
	// must be written if it exists.
	b.WriteString(s[start:])
	return b.String(), nil
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
