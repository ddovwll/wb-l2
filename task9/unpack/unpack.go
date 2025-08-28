package unpack

import (
	"errors"
	"strings"
	"unicode"
)

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var b strings.Builder
	runes := []rune(s)
	escape := false
	var prev rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if escape {
			b.WriteRune(r)
			prev = r
			escape = false
			continue
		}

		if r == '\\' {
			escape = true
			continue
		}

		if unicode.IsDigit(r) {
			if prev == 0 {
				return "", errors.New("invalid string")
			}
			count := int(r - '0')
			if count == 0 {
				tmp := []rune(b.String())
				b.Reset()
				b.WriteString(string(tmp[:len(tmp)-1]))
			} else {
				b.WriteString(strings.Repeat(string(prev), count-1))
			}
			continue
		}

		b.WriteRune(r)
		prev = r
	}

	if escape {
		return "", errors.New("invalid string")
	}

	return b.String(), nil
}
