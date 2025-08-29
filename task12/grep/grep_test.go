package grep

import (
	"strings"
	"testing"
)

func TestFixedMatch(t *testing.T) {
	lines := []string{"foo", "bar", "baz", "BaR"}
	opts := Options{Pattern: "bar", Fixed: true, IgnoreCase: true}

	result, err := Process(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"bar", "BaR"}
	if len(result) != len(want) || result[0] != want[0] {
		t.Errorf("got %v, want %v", result, want)
	}
}

func TestIgnoreCase(t *testing.T) {
	lines := []string{"Hello", "world"}
	opts := Options{Pattern: "hello", IgnoreCase: true}

	result, err := Process(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "Hello" {
		t.Errorf("case-insensitive match failed: got %v", result)
	}
}

func TestInvertMatch(t *testing.T) {
	lines := []string{"foo", "bar"}
	opts := Options{Pattern: "foo", Fixed: true, Invert: true}

	result, err := Process(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "bar" {
		t.Errorf("invert match failed: got %v", result)
	}
}

func TestCountOnly(t *testing.T) {
	lines := []string{"a", "b", "a", "c"}
	opts := Options{Pattern: "a", Fixed: true, CountOnly: true}

	result, err := Process(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result[0] != "2" {
		t.Errorf("countOnly failed: got %v, want 2", result[0])
	}
}

func TestContextBeforeAfter(t *testing.T) {
	lines := []string{
		"line1",
		"line2",
		"target",
		"line4",
		"line5",
	}
	opts := Options{Pattern: "target", Fixed: true, Before: 1, After: 1}

	result, err := Process(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"line2", "target", "line4"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("got %v, want %v", result, want)
		}
	}
}

func TestShowLineNumbers(t *testing.T) {
	lines := []string{"alpha", "beta", "gamma"}
	opts := Options{Pattern: "beta", Fixed: true, ShowNum: true}

	result, err := Process(lines, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result[0] != "2:beta" {
		t.Errorf("line number failed: got %v", result[0])
	}
}

func TestReadLines(t *testing.T) {
	data := "foo\nbar\nbaz\n"
	reader := strings.NewReader(data)

	lines, err := ReadLines(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"foo", "bar", "baz"}
	if len(lines) != len(want) {
		t.Fatalf("got %v, want %v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("got %v, want %v", lines[i], want[i])
		}
	}
}
