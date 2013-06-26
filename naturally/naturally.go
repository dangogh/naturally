// Implements "natural" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
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
        partsA := partition(p[a])
        partsB := partition(p[b])
        // if equal up to end of one, shorter of the two is less
        less := len(partsA) < len(partsB)
        for ii, ca := range partsA {
                if ii >= len(partsB) {
                        break
                }
                cb := partsB[ii]
                if ca != cb {
                        if ia, err := strconv.Atoi(ca); err != nil {
                                // non-numeric 
                                less = (ca < cb)
                        } else {
                                // numeric
                                ib, _ := strconv.Atoi(cb)
                                less = (ia < ib)
                        }
                        break
                }
        }
	return less
}
