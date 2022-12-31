package main

import (
	"bytes"
	"fmt"
	"os"
)

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

// Table represents a table.
type Table struct {
	// Head contains the name of each column.
	Head []string
	// Body contains each row of the table.
	Body []map[string]string
}

func main() {
	t, err := ReadTsv(os.Stdin, false)
	if err != nil {
		errExit(err)
	}
	switch os.Args[1] {
	case "template":
		f, err := os.Open(os.Args[2])
		if err != nil {
			errExit(err)
		}
		tmpl, err := ReadTemplate(f)
		if err != nil {
			errExit(err)
		}
		if err := tmpl.Execute(os.Stdout, t); err != nil {
			errExit(err)
		}
	case "tojson":
		fmt.Println(t.ToJson())
	case "totsv":
		tsv, err := t.ToTsv(false)
		if err != nil {
			errExit(err)
		}
		fmt.Print(tsv)
	default:
		errExit(fmt.Errorf("unrecognized command: %s", os.Args[1]))
	}
}

// errExit writes a formatted error message to stderr
// and exits with a non-zero exit code.
func errExit(e error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], e)
	os.Exit(1)
}

// scanPosixLines is like bufio.ScanLines
// but it adheres strictly to POSIX's definition of a line.
// Carriage returns are always considered part of a line.
// Data before EOF not terminated by a line feed is discarded.
func scanPosixLines(d []byte, atEOF bool) (int, []byte, error) {
	if i := bytes.IndexByte(d, '\n'); i >= 0 {
		return i + 1, d[:i], nil
	}
	// If atEOF, return nothing; otherwise, request more data.
	return 0, nil, nil
}
