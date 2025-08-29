package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"task10/sort"
	"unicode"
)

func main() {
	flag.Parse()

	var lines []string
	args := flag.Args()

	if len(args) > 0 {
		for _, fileName := range args {
			if err := readFileLines(fileName, &lines); err != nil {
				fmt.Fprintf(os.Stderr, "error while reading file %s: %v\n", fileName, err)
				os.Exit(1)
			}
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if *sort.IgnoreBlanks {
				line = strings.TrimRightFunc(line, unicode.IsSpace)
			}
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
			os.Exit(1)
		}
	}

	result := sort.ProcessLines(lines)
	for _, line := range result {
		fmt.Println(line)
	}
}

func readFileLines(fileName string, lines *[]string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error closing file %s: %v\n", fileName, err)
			os.Exit(1)
		}
	}(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if *sort.IgnoreBlanks {
			line = strings.TrimRightFunc(line, unicode.IsSpace)
		}
		*lines = append(*lines, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
