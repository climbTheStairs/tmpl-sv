package main

import (
	"bytes"
	"encoding/json"
	"strings"
)

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
