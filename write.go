package main

import (
	"fmt"
	"strings"
)

// ForbiddenCharacters contains all the characters
// forbidden in text/tab-separated-values.
var ForbiddenChars = []byte{'\n', '\t'}

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
				return "", fmt.Errorf(
					`row %d: column %d "%s": %v`,
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
			return field, fmt.Errorf(`forbidden character: %s`,
				Escape(sc))
		}
	}
	return field, nil
}

// ToJson returns table t as a JSON array.
func (t *Table) ToJson() string {
	arr := make([]string, len(t.Body))
	for i, row := range t.Body {
		obj := make([]string, len(t.Head))
		for i, k := range t.Head {
			obj[i] = escapeJson(k) + ":" + escapeJson(row[k])
		}
		arr[i] = "{" + strings.Join(obj, ",") + "}"
	}
	return "[" + strings.Join(arr, ",") + "]"
}
