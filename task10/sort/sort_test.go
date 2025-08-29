package sort

import (
	"flag"
	"reflect"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
	Column = flag.CommandLine.Int("k", 0, "")
	Numeric = flag.CommandLine.Bool("n", false, "")
	Reverse = flag.CommandLine.Bool("r", false, "")
	Unique = flag.CommandLine.Bool("u", false, "")
	MonthSort = flag.CommandLine.Bool("M", false, "")
	IgnoreBlanks = flag.CommandLine.Bool("b", false, "")
	CheckOnly = flag.CommandLine.Bool("c", false, "")
	HumanSort = flag.CommandLine.Bool("h", false, "")
}

func TestNumericSort(t *testing.T) {
	resetFlags()
	*Numeric = true

	lines := []string{"10", "2", "1", "30"}
	got := ProcessLines(lines)
	want := []string{"1", "2", "10", "30"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Numeric sort: got %v, want %v", got, want)
	}
}

func TestReverseSort(t *testing.T) {
	resetFlags()
	*Reverse = true

	lines := []string{"a", "c", "b"}
	got := ProcessLines(lines)
	want := []string{"c", "b", "a"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Reverse sort: got %v, want %v", got, want)
	}
}

func TestColumnSort(t *testing.T) {
	resetFlags()
	*Column = 2
	*Numeric = true

	lines := []string{
		"foo\t3",
		"baz\t1",
		"bar\t2",
	}
	got := ProcessLines(lines)
	want := []string{
		"baz\t1",
		"bar\t2",
		"foo\t3",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Column sort: got %v, want %v", got, want)
	}
}

func TestMonthSort(t *testing.T) {
	resetFlags()
	*MonthSort = true

	lines := []string{"Mar", "Jan", "Feb"}
	got := ProcessLines(lines)
	want := []string{"Jan", "Feb", "Mar"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Month sort: got %v, want %v", got, want)
	}
}

func TestHumanSort(t *testing.T) {
	resetFlags()
	*HumanSort = true

	lines := []string{"10K", "2M", "512"}
	got := ProcessLines(lines)
	want := []string{"512", "10K", "2M"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Human sort: got %v, want %v", got, want)
	}
}

func TestUnique(t *testing.T) {
	resetFlags()
	*Unique = true

	lines := []string{"a", "b", "a", "c", "b"}
	got := ProcessLines(lines)
	want := []string{"a", "b", "c"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unique: got %v, want %v", got, want)
	}
}

func TestIsSorted(t *testing.T) {
	lines := []string{"1", "2", "3"}
	if !isSorted(lines, less) {
		t.Errorf("expected sorted, but got not sorted")
	}
	lines = []string{"3", "1", "2"}
	if isSorted(lines, less) {
		t.Errorf("expected not sorted, but got sorted")
	}
}
