package cut

import (
	"bytes"
	"strings"
	"testing"
)

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		opts     Options
		expected string
		ok       bool
	}{
		{
			name:     "simple one field",
			line:     "a:b:c",
			opts:     Options{Fields: []int{2}, Delimiter: ":", Separated: false},
			expected: "b",
			ok:       true,
		},
		{
			name:     "multiple fields",
			line:     "a:b:c",
			opts:     Options{Fields: []int{1, 3}, Delimiter: ":", Separated: false},
			expected: "a:c",
			ok:       true,
		},
		{
			name:     "field out of range",
			line:     "a:b",
			opts:     Options{Fields: []int{3}, Delimiter: ":", Separated: false},
			expected: "",
			ok:       false,
		},
		{
			name:     "Separated=true, no delimiter",
			line:     "abc",
			opts:     Options{Fields: []int{1}, Delimiter: ":", Separated: true},
			expected: "",
			ok:       false,
		},
		{
			name:     "Separated=false, no delimiter",
			line:     "abc",
			opts:     Options{Fields: []int{1}, Delimiter: ":", Separated: false},
			expected: "abc",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, ok := ProcessLine(tt.line, tt.opts)
			if out != tt.expected || ok != tt.ok {
				t.Errorf("got (%q,%v), expected (%q,%v)", out, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestRun(t *testing.T) {
	input := "a:b:c\nd:e:f\n"
	opts := Options{Fields: []int{2}, Delimiter: ":", Separated: false}
	var output bytes.Buffer

	err := Run(strings.NewReader(input), &output, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "b\ne\n"
	if output.String() != expected {
		t.Errorf("got %q, expected %q", output.String(), expected)
	}
}

func TestParseFields(t *testing.T) {
	tests := []struct {
		spec     string
		expected []int
		hasError bool
	}{
		{"1", []int{1}, false},
		{"2,4", []int{2, 4}, false},
		{"1-3", []int{1, 2, 3}, false},
		{"2-2", []int{2}, false},
		{"1,3-5", []int{1, 3, 4, 5}, false},
		{"0", nil, true},
		{"a", nil, true},
		{"2-1", nil, true},
		{"1-", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.spec, func(t *testing.T) {
			fields, err := ParseFields(tt.spec)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error=%v, got %v", tt.hasError, err)
			}
			if !tt.hasError {
				if len(fields) != len(tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, fields)
				}
				for i := range fields {
					if fields[i] != tt.expected[i] {
						t.Errorf("expected %v, got %v", tt.expected, fields)
					}
				}
			}
		})
	}
}
