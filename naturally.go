// Package naturally implements "natural" sort ordering: the alphabetic
// portions of a string are compared alphabetically and the numeric portions
// are compared numerically, so that "A2" sorts before "A11".
//
// Ordering details worth knowing:
//
//   - Numeric segments are compared by value, with no limit on the number of
//     digits, so "99999999999999999999" sorts before "100000000000000000000".
//   - When two numeric segments have the same value but different lengths, the
//     shorter (less zero-padded) one sorts first: "A1" < "A01" < "A001".
//   - A string containing a digit sorts before one that does not, so
//     "b2" < "Banana". See the package README for the rationale.
//   - Segments of non-ASCII digits (for example Arabic-Indic "٣") are compared
//     as text rather than by value.
package naturally

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// StringSlice attaches the methods of sort.Interface to []string, sorting in
// natural order (case-sensitive).
type StringSlice sort.StringSlice

// Len is the number of elements in the collection.
func (p StringSlice) Len() int { return len(p) }

// Swap swaps the elements with indexes a and b.
func (p StringSlice) Swap(a, b int) { p[a], p[b] = p[b], p[a] }

func isNonDigit(ch rune) bool {
	return !unicode.IsDigit(ch)
}

// Less reports whether element a should sort before element b.
func (p StringSlice) Less(a, b int) bool {
	return less(p[a], p[b], false)
}

// CIStringSlice attaches the methods of sort.Interface to []string, sorting in
// natural order while ignoring case.
type CIStringSlice sort.StringSlice

// Len is the number of elements in the collection.
func (p CIStringSlice) Len() int { return len(p) }

// Swap swaps the elements with indexes a and b.
func (p CIStringSlice) Swap(a, b int) { p[a], p[b] = p[b], p[a] }

// Less reports whether element a should sort before element b, ignoring case.
func (p CIStringSlice) Less(a, b int) bool {
	return less(p[a], p[b], true)
}

// Sort sorts a string slice using natural ordering (case-sensitive).
func Sort(s []string) {
	sort.Sort(StringSlice(s))
}

// SortCI sorts a string slice using natural ordering (case-insensitive).
func SortCI(s []string) {
	sort.Sort(CIStringSlice(s))
}

func less(strA, strB string, caseInsensitive bool) bool {
	if strA == strB {
		return false
	}
	for {
		// get chars up to 1st digit
		posA := strings.IndexFunc(strA, unicode.IsDigit)
		posB := strings.IndexFunc(strB, unicode.IsDigit)

		if posA == -1 {
			// no digits in A
			if posB == -1 {
				// or B -- straight string compare
				if caseInsensitive {
					return compareFold(strA, strB) < 0
				}
				return strA < strB
			}
			return false // B is Less
		} else if posB == -1 {
			return true // A is Less
		}
		subA, subB := strA[:posA], strB[:posB]
		if caseInsensitive {
			if !strings.EqualFold(subA, subB) {
				return compareFold(subA, subB) < 0
			}
		} else if subA != subB {
			return subA < subB
		}
		strA, strB = strA[posA:], strB[posB:]

		// get chars up to 1st non-digit
		posA = strings.IndexFunc(strA, isNonDigit)
		posB = strings.IndexFunc(strB, isNonDigit)
		if posA == -1 {
			posA = len(strA)
		}
		if posB == -1 {
			posB = len(strB)
		}

		// grab numeric part of each
		numA, numB := strA[:posA], strB[:posB]
		if isASCIIDigits(numA) && isASCIIDigits(numB) {
			if cmp := compareNumeric(numA, numB); cmp != 0 {
				return cmp < 0
			}
		} else if numA != numB {
			// non-ASCII digits carry no usable value: compare them as text
			return numA < numB
		}
		if posA != posB {
			return posA < posB
		}
		if posA >= len(strA) || posB >= len(strB) {
			return len(strA) < len(strB)
		}
		strA, strB = strA[posA:], strB[posB:]
	}
}

func isASCIIDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// compareNumeric compares two non-empty runs of ASCII digits by value,
// returning -1, 0, or 1. It handles arbitrarily long runs, so it does not
// overflow the way strconv.Atoi does. Leading zeros do not affect the value;
// callers break that tie on digit-run length.
func compareNumeric(a, b string) int {
	a = strings.TrimLeft(a, "0")
	b = strings.TrimLeft(b, "0")
	if len(a) != len(b) {
		if len(a) < len(b) {
			return -1
		}
		return 1
	}
	return strings.Compare(a, b)
}

// compareFold compares a and b by their lowercased runes, returning -1, 0, or
// 1. It matches the ordering of comparing strings.ToLower(a) to
// strings.ToLower(b) without allocating.
func compareFold(a, b string) int {
	for a != "" && b != "" {
		ra, sizeA := utf8.DecodeRuneInString(a)
		rb, sizeB := utf8.DecodeRuneInString(b)
		la, lb := unicode.ToLower(ra), unicode.ToLower(rb)
		if la != lb {
			if la < lb {
				return -1
			}
			return 1
		}
		a, b = a[sizeA:], b[sizeB:]
	}
	switch {
	case a == "" && b == "":
		return 0
	case a == "":
		return -1
	default:
		return 1
	}
}
