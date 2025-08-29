package cut

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Options struct {
	Fields    []int
	Delimiter string
	Separated bool
}

func ProcessLine(line string, opts Options) (string, bool) {
	if opts.Separated && !strings.Contains(line, opts.Delimiter) {
		return "", false
	}

	parts := strings.Split(line, opts.Delimiter)

	var selected []string
	for _, f := range opts.Fields {
		if f-1 >= 0 && f-1 < len(parts) {
			selected = append(selected, parts[f-1])
		}
	}
	if len(selected) == 0 {
		return "", false
	}

	return strings.Join(selected, opts.Delimiter), true
}

func Run(r io.Reader, w io.Writer, opts Options) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if out, ok := ProcessLine(line, opts); ok {
			_, err := fmt.Fprintln(w, out)
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func ParseFields(spec string) ([]int, error) {
	var fields []int
	parts := strings.Split(spec, ",")
	for _, p := range parts {
		if strings.Contains(p, "-") {
			bounds := strings.Split(p, "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start <= 0 || end < start {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			num, err := strconv.Atoi(p)
			if err != nil || num <= 0 {
				return nil, fmt.Errorf("invalid field number: %s", p)
			}
			fields = append(fields, num)
		}
	}
	return fields, nil
}
