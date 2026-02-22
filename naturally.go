// Implements "natural" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Naturally implements sort.Interface
type StringSlice sort.StringSlice

func (p StringSlice) Len() int      { return len(p) }
func (p StringSlice) Swap(a, b int) { p[a], p[b] = p[b], p[a] }

func isNonDigit(ch rune) bool {
	return !unicode.IsDigit(ch)
}

func (p StringSlice) Less(a, b int) bool {
	strA := p[a][:]
	strB := p[b][:]
	//fmt.Println(strA, " <=> ", strB)

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
				return strA < strB
			}
			return false // B is Less
		} else if posB == -1 {
			return true // A is Less
		}
		subA, subB := strA[:posA], strB[:posB]
		if subA != subB {
			return subA < subB
		}
		strA, strB = strA[posA:], strB[posB:]

		// get chars up to 1st non-digit
		posA = strings.IndexFunc(strA, isNonDigit)
		posB = strings.IndexFunc(strB, isNonDigit)
		if posA == -1 {
			// no non-digits in A - allow numeric compare
			//fmt.Println(posA, " pos in ", strA)
			posA = len(strA)
		}
		if posB == -1 {
			// no non-digits in B - allow numeric compare
			posB = len(strB)
		}

		// grab numeric part of each
		valA, errA := strconv.Atoi(strA[:posA])
		valB, errB := strconv.Atoi(strB[:posB])
		if errA != nil || errB != nil {
			// fall back to string comparison for unparseable digits (e.g. non-ASCII)
			return strA[:posA] < strB[:posB]
		}
		if valA != valB {
			return valA < valB
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
