package main

import (
	"strings"
)

// ToTsv returns table t in TSV format.
func (t *Table) ToTsv() string {
	lines := make([]string, 0, len(t.Body)+1)
	for i := 0; i < len(t.Body)+1; i++ {
		line := make([]string, 0, len(t.Head))
		for _, v := range t.Head {
			if i > 0 {
				v = t.Body[i-1][v]
			}
			line = append(line, Escape(v))
		}
		lines = append(lines, strings.Join(line, "\t"))
	}
	return strings.Join(lines, "\n") + "\n"
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
