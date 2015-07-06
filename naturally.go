// Implements "natural" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
        "strings"
	"strconv"
)

// Naturally implements sort.Interface
type StringSlice []string

func (p StringSlice) Len() int { return len(p) }
func (p StringSlice) Swap(a, b int) { p[a], p[b] = p[b], p[a] }

func partition(s string) (parts []string ) {
        isnumeric := false
        last := 0
        for ii, c := range s {
                isdigit := (c >= '0' && c <= '9')
                if ii == last {
                        isnumeric = isdigit
                        continue
                }
                if isnumeric != isdigit {
                        parts = append(parts, s[last:ii])
                        isnumeric = isdigit
                        last = ii
                }
        }
        ii := len(s)
        parts = append(parts, s[last:ii])
        return parts
}

func (p StringSlice) Less(a, b int) bool {
	// part string -- numeric and non
        sA,sB := p[a][0:], p[b][0:]
        digits := "0123456789"
        var idxA, idxB int
        for {
                idxA = strings.IndexAny(sA, digits)
                idxB = strings.IndexAny(sB, digits)
                switch {
                case idxA == -1:

                case idxB == -1:
                        return false
                case sA[0:idxA] == sA[0:idxB]:
                        sA, sB = sA[idxA+1:], sB[idxB+1:]
                        continue
                case idxA == 0:
                        return strconv.Atoi(sA) < strconv.Atoi(sB)
                default:
                        return sA < sB
                }
        }
}
