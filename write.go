package main

import (
	"fmt"
	"strings"
)

// ToTsv returns table t in TSV format.
func (t *Table) ToTsv() string {
	lines := make([]string, 0, len(t.Body)+1)

	line := make([]string, 0, len(t.Head))
	for _, v := range t.Head {
		line = append(line, Escape(v))
	}
	lines = append(lines, strings.Join(line, "\t"))

	for _, row := range t.Body {
		line := make([]string, 0, len(t.Head))
		for _, k := range t.Head {
			line = append(line, Escape(row[k]))
		}
		lines = append(lines, strings.Join(line, "\t"))
	}

	return strings.Join(lines, "\n") + "\n"
}

// ToJson returns table t as a JSON array.
// Escaping special JSON characters is yet to be implemented!!!
func (t *Table) ToJson() string {
	arr := make([]string, len(t.Body))
	for i, row := range t.Body {
		obj := make([]string, len(t.Head))
		for i, k := range t.Head {
			obj[i] = fmt.Sprintf(`"%s":"%s"`, k, row[k])
		}
		arr[i] = "{" + strings.Join(obj, ",") + "}"
	}
	return "[" + strings.Join(arr, ",") + "]"
}
