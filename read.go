package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	table, err := readTsv(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if err := applyTmpl(os.Args[1], table, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

// scanPosixLines is like bufio.ScanLines
// but it adheres strictly to POSIX's definition of a line.
// Carriage returns are always considered part of a line.
// Data before EOF not terminated by a line feed is discarded.
func scanPosixLines(d []byte, atEOF bool) (int, []byte, error) {
	if i := bytes.IndexByte(d, '\n'); i >= 0 {
		return i+1, d[:i], nil
	}
	// If atEOF, return nothing; otherwise, request more data.
	return 0, nil, nil
}

func readTsv(f io.Reader) ([]map[string]string, error) {
	scanner := bufio.NewScanner(f)
	scanner.Split(scanPosixLines)

	if ok := scanner.Scan(); !ok {
		return nil, errors.New("readTsv: cannot read first line")
	}
	head := strings.Split(scanner.Text(), "\t")

	table := make([]map[string]string, 0)
	for scanner.Scan() {
		row := make(map[string]string)
		for i, v := range strings.Split(scanner.Text(), "\t") {
			row[head[i]] = v
		}
		table = append(table, row)
	}

	return table, nil
}
