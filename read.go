package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	table, err := readTsv(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	if err := applyTmpl(os.Args[1], table, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
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

func readTsv(f io.Reader) ([]map[string]string, error) {
	scanner := bufio.NewScanner(f)
	scanner.Split(scanPosixLines)

	if ok := scanner.Scan(); !ok {
		return nil, errors.New("cannot read initial line")
	}
	head := strings.Split(scanner.Text(), "\t")

	table := make([]map[string]string, 0)
	for lnum := 2; scanner.Scan(); lnum++ {
		cols := strings.Split(scanner.Text(), "\t")
		if len(cols) != len(head) {
			return table, fmt.Errorf("line %d contains "+
				"invalid number of columns", lnum)
		}
		row := make(map[string]string)
		for i, colName := range head {
			row[colName] = cols[i]
		}
		table = append(table, row)
	}

	return table, nil
}
