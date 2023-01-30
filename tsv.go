package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ForbiddenCharacters contains all the characters
// forbidden in text/tab-separated-values.
var ForbiddenChars = []byte{'\n', '\t'}

// Escapes is a map of the second character
// of each two-character escape sequence
// (i.e. the character that isn't a backslash ('\\'))
// to the special character
// represented by the whole escape sequence.
var Escapes = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\\': '\\',
}

var escaper *strings.Replacer

func init() {
	oldnew := make([]string, 0, len(Escapes)*2)
	for escaped, unescaped := range Escapes {
		oldnew = append(oldnew, string(unescaped), "\\"+string(escaped))
	}
	escaper = strings.NewReplacer(oldnew...)
}

// ReadTsv reads TSV data from r
// and creates and returns a new table.
func ReadTsv(r io.Reader, esc bool) (*Table, error) {
	var err error
	t := &Table{}

	scanner := bufio.NewScanner(r)
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

// ToTsv returns table t in TSV format.
// If esc is true, ToTsv will call Escape on each field.
// ToTsv will return a non-nil error
// if and only if esc is false
// and a field contains characters in ForbiddenChars.
func (t *Table) ToTsv(esc bool) (string, error) {
	var err error
	// len(t.Body) is added 1 because t.Head
	// becomes the first line.
	lines := make([]string, len(t.Body)+1)
	for i := 0; i < len(t.Body)+1; i++ {
		line := make([]string, len(t.Head))
		for j, k := range t.Head {
			v := k
			if i > 0 {
				v = t.Body[i-1][k]
			}
			v, err = escapeOrValidateField(v, esc)
			if err != nil {
				return "", fmt.Errorf(`row %d: column %d "%s": %v`,
					i, j+1, k, err)
			}
			line[j] = v
		}
		lines[i] = strings.Join(line, "\t")
	}
	return strings.Join(lines, "\n") + "\n", nil
}

// escapeOrValidateField is a helper function for Table.ToTsv.
// escapeOrValidateField escapes field, if necessary;
// otherwise, it checks that field does not contain
// characters that must be escaped.
func escapeOrValidateField(field string, esc bool) (string, error) {
	if esc {
		return Escape(field), nil
	}
	for _, c := range ForbiddenChars {
		sc := string(c)
		if strings.Contains(field, sc) {
			return field, fmt.Errorf(`forbidden character: %s`, Escape(sc))
		}
	}
	return field, nil
}

// Escape returns a copy of string s
// with each special character disallowed in TSV
// replaced by its two-character escape sequnce.
// See Escapes for specific characters that are escaped
// and their escape sequences.
func Escape(s string) string {
	return escaper.Replace(s)
}
