package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Options struct {
	After      int
	Before     int
	CountOnly  bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	ShowNum    bool
	Pattern    string
}

func matchFunc(pattern string, opts Options) (func(string) bool, error) {
	if opts.Fixed {
		if opts.IgnoreCase {
			pattern = strings.ToLower(pattern)
			return func(s string) bool {
				return strings.Contains(strings.ToLower(s), pattern)
			}, nil
		}
		return func(s string) bool {
			return strings.Contains(s, pattern)
		}, nil
	}

	flags := ""
	if opts.IgnoreCase {
		flags = "(?i)"
	}
	re, err := regexp.Compile(flags + pattern)
	if err != nil {
		return nil, err
	}
	return func(s string) bool {
		return re.MatchString(s)
	}, nil
}

func ReadLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func Process(lines []string, opts Options) ([]string, error) {
	isMatch, err := matchFunc(opts.Pattern, opts)
	if err != nil {
		return nil, err
	}

	matches := make([]bool, len(lines))
	for i, line := range lines {
		ok := isMatch(line)
		if opts.Invert {
			ok = !ok
		}
		matches[i] = ok
	}

	if opts.CountOnly {
		cnt := 0
		for _, m := range matches {
			if m {
				cnt++
			}
		}
		return []string{fmt.Sprintf("%d", cnt)}, nil
	}

	printed := make([]bool, len(lines))
	var result []string
	for i, m := range matches {
		if !m {
			continue
		}
		start := i - opts.Before
		if start < 0 {
			start = 0
		}
		end := i + opts.After
		if end >= len(lines) {
			end = len(lines) - 1
		}
		for j := start; j <= end; j++ {
			if !printed[j] {
				if opts.ShowNum {
					result = append(result, fmt.Sprintf("%d:%s", j+1, lines[j]))
				} else {
					result = append(result, lines[j])
				}
				printed[j] = true
			}
		}
	}
	return result, nil
}
