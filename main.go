package main

import (
	"bytes"
	"fmt"
	"os"
)

type Table struct {
	Head []string
	Body []map[string]string
}

func main() {
	t, err := readTsv(os.Stdin)
	if err != nil {
		errExit(err)
	}
	if esc := true; esc {
		if err := escapeTable(&t); err != nil {
			errExit(err)
		}
	}
	if err := applyTmpl(os.Args[1], &t, os.Stdout); err != nil {
		errExit(err)
	}
}

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
