package sort

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	Column       = flag.Int("k", 0, "column sort")
	Numeric      = flag.Bool("n", false, "numeric sort")
	Reverse      = flag.Bool("r", false, "reverse order")
	Unique       = flag.Bool("u", false, "unique lines")
	MonthSort    = flag.Bool("M", false, "month sort")
	IgnoreBlanks = flag.Bool("b", false, "ignore trailing blanks")
	CheckOnly    = flag.Bool("c", false, "check if sorted")
	HumanSort    = flag.Bool("h", false, "human-readable (10K, 2M)")
)

var months = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
	"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
	"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

func ProcessLines(lines []string) []string {
	if *CheckOnly {
		if isSorted(lines, less) {
			fmt.Println("input is sorted")
		} else {
			fmt.Println("input is not sorted")
		}
		os.Exit(0)
	}

	sort.Slice(lines, func(i, j int) bool {
		return less(lines[i], lines[j])
	})

	if *Unique {
		lines = uniq(lines)
	}

	return lines
}

func less(a, b string) bool {
	ka := getKey(a)
	kb := getKey(b)

	var res int
	switch {
	case *Numeric:
		af, _ := strconv.ParseFloat(ka, 64)
		bf, _ := strconv.ParseFloat(kb, 64)
		if af < bf {
			res = -1
		} else if af > bf {
			res = 1
		}
	case *MonthSort:
		res = compareInt(months[ka], months[kb])
	case *HumanSort:
		af := parseHuman(ka)
		bf := parseHuman(kb)
		res = compareInt(af, bf)
	default:
		res = strings.Compare(ka, kb)
	}

	if *Reverse {
		return res > 0
	}
	return res < 0
}

func getKey(line string) string {
	cols := strings.Split(line, "\t")
	if *Column > 0 && *Column <= len(cols) {
		return cols[*Column-1]
	}
	return line
}

func compareInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func parseHuman(s string) int {
	if s == "" {
		return 0
	}
	n := len(s)
	unit := s[n-1]
	val, _ := strconv.ParseFloat(s[:n-1], 64)
	switch unit {
	case 'K', 'k':
		return int(val * 1024)
	case 'M', 'm':
		return int(val * 1024 * 1024)
	case 'G', 'g':
		return int(val * 1024 * 1024 * 1024)
	default:
		v, _ := strconv.Atoi(s)
		return v
	}
}

func uniq(lines []string) []string {
	var result []string
	seen := make(map[string]struct{})
	for _, l := range lines {
		if _, ok := seen[l]; !ok {
			result = append(result, l)
			seen[l] = struct{}{}
		}
	}
	return result
}

func isSorted(lines []string, cmp func(a, b string) bool) bool {
	for i := 1; i < len(lines); i++ {
		if cmp(lines[i], lines[i-1]) {
			return false
		}
	}
	return true
}
