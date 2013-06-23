// Implements "Naturally" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
        "fmt"
	//"strconv"
)

// Naturally implements sort.Interface by providing Less and
// using Len and Swap  methods of embedded []string
type Naturally struct {
	Val []string
}

func (p Naturally) Len() int {
	return len(p.Val)
}

func (p Naturally) Swap(a, b int) {
	p.Val[a], p.Val[b] = p.Val[b], p.Val[a]
}

func trace(isnumeric bool, val string) {
        t := "alpha"
        if isnumeric {
                t = "num"
        }
        fmt.Printf("%s %v\n", t, val)
}

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
                        trace(isnumeric, s[last:ii])
                        parts = append(parts, s[last:ii])
                        isnumeric = isdigit
                        last = ii
                }
        }
        ii := len(s)
        trace(isnumeric, s[last:ii])
        parts = append(parts, s[last:ii])
        return parts
}

func (p Naturally) Less(a, b int) bool {
	// part string -- numeric and non
        partsA := partition(p.Val[a])
        partsB := partition(p.Val[b])
        less := true
        for ii, ca := range partsA {
                if ii >= len(partsB) {
                        break
                }
                cb := partsB[ii]
                if ca != cb {
                        less = (ca < cb)
                        break
                }
        }
	return less
}
