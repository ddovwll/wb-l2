package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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
				log.Fatal(err)
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
			log.Fatal(err)
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
			log.Fatal(err)
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
