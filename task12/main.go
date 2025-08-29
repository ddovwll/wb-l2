package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"task12/grep"
)

func main() {
	after := flag.Int("A", 0, "Print N lines after match")
	before := flag.Int("B", 0, "Print N lines before match")
	context := flag.Int("C", 0, "Print N lines of context around match")
	countOnly := flag.Bool("c", false, "Print only count of matching lines")
	ignoreCase := flag.Bool("i", false, "Ignore case distinctions")
	invert := flag.Bool("v", false, "Invert match")
	fixed := flag.Bool("F", false, "Fixed string match")
	showNum := flag.Bool("n", false, "Show line numbers")

	flag.Parse()

	if flag.NArg() < 1 {
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	var filename string
	if flag.NArg() > 1 {
		filename = flag.Arg(1)
	}

	if *context > 0 {
		*after, *before = *context, *context
	}

	var input *os.File
	var err error
	if filename != "" {
		input, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer func(input *os.File) {
			err := input.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(input)
	} else {
		input = os.Stdin
	}

	lines, err := grep.ReadLines(input)
	if err != nil {
		log.Fatal(err)
	}

	opts := grep.Options{
		After:      *after,
		Before:     *before,
		CountOnly:  *countOnly,
		IgnoreCase: *ignoreCase,
		Invert:     *invert,
		Fixed:      *fixed,
		ShowNum:    *showNum,
		Pattern:    pattern,
	}

	results, err := grep.Process(lines, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range results {
		fmt.Println(line)
	}
}
