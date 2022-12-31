package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ForbiddenCharacters contains all the characters
// forbidden in text/tab-separated-values.
var ForbiddenChars = []byte{'\n', '\t'}

var escaper *strings.Replacer

func init() {
	oldnew := make([]string, 0, len(Escapes)*2)
	for escaped, unescaped := range Escapes {
		oldnew = append(oldnew, string(unescaped), "\\"+string(escaped))
	}
	escaper = strings.NewReplacer(oldnew...)
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

// escapeJson wraps goodJsonMarshal, panicking on errors.
func escapeJson(s string) string {
	b, err := goodJsonMarshal(s)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// goodJsonMarshal is like json.Marshal but good.
// goodJsonMarshal is identical to json.Marshal
// but without its annoying and unasked-for escaping of characters
// that unnecessarily attempts to make the output HTML-safe.
func goodJsonMarshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	// Remove extra newline stupidly added by json.Encoder.Encode.
	b := buf.Bytes()[:buf.Len()-1]
	return b, err
}
