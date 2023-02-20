package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ReadJson reads JSON data from r
// and creates and returns a new table.
// I know this does not work for all inputs
// and I don't know if this works at all.
func ReadJson(r io.Reader) (*Table, error) {
	var t *Table

	var d []map[string]any
	dec := json.NewDecoder(r)
	if err := dec.Decode(&d); err != nil {
		return t, err
	}
	if len(d) == 0 {
		return t, fmt.Errorf("no data")
	}

	t.Head = make([]string, 0, len(d[0]))
	for k := range d[0] {
		t.Head = append(t.Head, k)
	}
	sort.Strings(t.Head)

	for i, fields := range d {
		i += 1
		if len(fields) != len(t.Head) {
			return t, fmt.Errorf("row %d: invalid number of fields", i)
		}
		row := make(map[string]string)
		for _, k := range t.Head {
			v, ok := fields[k]
			if !ok {
				return t, fmt.Errorf("%d...", i+1)
			}
			switch s := v.(type) {
			case string:
				row[k] = s
			case bool:
				row[k] = "0"
				if s {
					row[k] = "1"
				}
			case json.Number:
				row[k] = string(s)
			default:
				if s != nil {
					return t, fmt.Errorf("...")
				}
				row[k] = ""
			}
		}
		t.Body = append(t.Body, row)
	}
	return t, nil
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

// escapeJson wraps goodJsonMarshal.
func escapeJson(s string) string {
	b, err := goodJsonMarshal(s)
	if err != nil {
		// goodJsonMarshal should never error with a string;
		// therefore, it is safe to do whatever with err
		// because this will never occur.
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
	// Remove extra newline stupidly added by json.Encoder.Encode
	// which is different for some reason from json.Marshal
	// and the devs [REFUSE to fix]!
	//
	// [REFUSE to fix]: https://github.com/golang/go/issues/37083
	b := buf.Bytes()[:buf.Len()-1]
	return b, err
}
