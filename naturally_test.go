package naturally_test

import (
	"sort"
	"testing"

	"github.com/dangogh/naturally"
)

func TestLess(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		// identical strings
		{"equal strings", "abc", "abc", false},
		{"equal with digits", "A1B2", "A1B2", false},

		// pure alphabetic (no digits at all)
		{"alpha a<b", "abc", "def", true},
		{"alpha a>b", "def", "abc", false},
		{"alpha prefix shorter", "A", "AA", true},
		{"alpha prefix longer", "AA", "A", false},

		// no digits in A, digits in B -> A is not less (returns false)
		{"no digits A, digits B", "abc", "a1", false},

		// digits in A, no digits in B -> A is less (returns true)
		{"digits A, no digits B", "a1", "abc", true},

		// different text prefixes before first digit
		{"diff prefix a<b", "A1", "B1", true},
		{"diff prefix a>b", "B1", "A1", false},

		// same prefix, different numbers
		{"same prefix smaller num", "A1", "A2", true},
		{"same prefix larger num", "A2", "A1", false},
		{"same prefix multi-digit", "A2", "A11", true},
		{"same prefix multi-digit rev", "A11", "A2", false},

		// leading zeros -- same numeric value but different string lengths
		{"leading zeros 01 vs 1", "A01", "A1", false},
		{"leading zeros 1 vs 01", "A1", "A01", true},
		{"leading zeros 001 vs 01", "A001", "A01", false},
		{"leading zeros 01 vs 001", "A01", "A001", true},

		// number at start of string (empty text prefix)
		{"num at start", "1abc", "2abc", true},
		{"num at start rev", "2abc", "1abc", false},
		{"num at start equal", "1abc", "1def", true},
		{"num at start equal rev", "1def", "1abc", false},

		// multiple segments -- text+num+text+num
		{"multi segment same first", "A1B2", "A1B3", true},
		{"multi segment same first rev", "A1B3", "A1B2", false},
		{"multi segment diff second text", "A1A1", "A1B1", true},
		{"multi segment diff second text rev", "A1B1", "A1A1", false},

		// string ends at numeric part (no trailing text)
		{"ends at number a<b", "A1", "A2", true},
		{"ends at number equal num", "A1", "A1", false},

		// one string ends in number, other continues
		{"a ends at num b continues", "A1", "A1B", true},
		{"a continues b ends at num", "A1B", "A1", false}, // A1B > A1: longer string with extra content sorts after

		// purely numeric strings
		{"pure num 1<2", "1", "2", true},
		{"pure num 2>1", "2", "1", false},
		{"pure num 9<10", "9", "10", true},
		{"pure num 10>9", "10", "9", false},
		{"pure num same value diff zeros", "00", "0", false},
		{"pure num same value diff zeros rev", "0", "00", true},

		// empty strings
		{"both empty", "", "", false},
		{"a empty b not", "", "a", true},
		{"a not b empty", "a", "", false},

		// non-ASCII digits (fallback to string comparison)
		{"arabic-indic digits", "A٣", "A٤", true},
		{"arabic-indic vs ascii", "A٣", "A3", false},

		// mixed case
		{"case sensitive lower<upper", "a1", "B1", false},
		{"case sensitive upper<lower", "B1", "a1", true},

		// large numbers
		{"large numbers", "file999", "file1000", true},
		{"large numbers rev", "file1000", "file999", false},

		// adjacent numeric segments with different text separators
		{"adj segments", "X1Y2", "X1Z1", true},
		{"adj segments rev", "X1Z1", "X1Y2", false},

		// realistic filenames
		{"file names", "photo9.jpg", "photo10.jpg", true},
		{"file names rev", "photo10.jpg", "photo9.jpg", false},
		{"file names same num", "photo10.jpg", "photo10.png", true},
		{"versioned", "v1.2.3", "v1.10.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := naturally.StringSlice{tt.a, tt.b}
			got := s.Less(0, 1)
			if got != tt.want {
				t.Errorf("Less(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name string
		data naturally.StringSlice
		want int
	}{
		{"empty", naturally.StringSlice{}, 0},
		{"one", naturally.StringSlice{"a"}, 1},
		{"three", naturally.StringSlice{"a", "b", "c"}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Len(); got != tt.want {
				t.Errorf("Len() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	s := naturally.StringSlice{"x", "y", "z"}
	s.Swap(0, 2)
	if s[0] != "z" || s[2] != "x" {
		t.Errorf("Swap(0,2) got %v, want [z y x]", s)
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			"alpha with numbers",
			[]string{"A1", "A0", "A21", "A11", "A111", "A21", "A2"},
			[]string{"A0", "A1", "A2", "A11", "A21", "A21", "A111"},
		},
		{
			"numeric with leading zeros",
			[]string{"001", "2", "xyzzy", "30", "22", "0", "00", "3", "1"},
			[]string{"0", "00", "1", "001", "2", "3", "22", "30", "xyzzy"},
		},
		{
			"multi-segment",
			[]string{"A1BA1", "A11AA1", "A2AB0", "B1AA1", "A1AA1"},
			[]string{"A1AA1", "A1BA1", "A2AB0", "A11AA1", "B1AA1"},
		},
		{
			"pure alpha",
			[]string{"ZZ", "XX", "AAA", "A", "AA", "Z", "XX"},
			[]string{"A", "AA", "AAA", "XX", "XX", "Z", "ZZ"},
		},
		{
			"empty slice",
			[]string{},
			[]string{},
		},
		{
			"single element",
			[]string{"only"},
			[]string{"only"},
		},
		{
			"already sorted",
			[]string{"a1", "a2", "a3"},
			[]string{"a1", "a2", "a3"},
		},
		{
			"reverse sorted",
			[]string{"a3", "a2", "a1"},
			[]string{"a1", "a2", "a3"},
		},
		{
			"realistic filenames",
			[]string{"img12.png", "img2.png", "img1.png", "img10.png", "img3.png"},
			[]string{"img1.png", "img2.png", "img3.png", "img10.png", "img12.png"},
		},
		{
			"mixed empty and non-empty",
			[]string{"b", "", "a"},
			[]string{"", "a", "b"},
		},
		{
			"duplicates",
			[]string{"a", "a", "a"},
			[]string{"a", "a", "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]string, len(tt.input))
			copy(data, tt.input)
			sort.Sort(naturally.StringSlice(data))

			if len(data) != len(tt.expected) {
				t.Fatalf("got len %d, want %d", len(data), len(tt.expected))
			}
			for i := range data {
				if data[i] != tt.expected[i] {
					t.Errorf("index %d: got %q, want %q (full result: %v)", i, data[i], tt.expected[i], data)
					break
				}
			}
		})
	}
}

func BenchmarkSortNaturally(b *testing.B) {
	src := []string{"A1BA1", "A11AA1", "A2AB0", "B1AA1", "A1AA1"}
	for i := 0; i < b.N; i++ {
		data := make([]string, len(src))
		copy(data, src)
		sort.Sort(naturally.StringSlice(data))
	}
}

func BenchmarkSortStdlib(b *testing.B) {
	src := []string{"A1BA1", "A11AA1", "A2AB0", "B1AA1", "A1AA1"}
	for i := 0; i < b.N; i++ {
		data := make([]string, len(src))
		copy(data, src)
		sort.Sort(sort.StringSlice(data))
	}
}
