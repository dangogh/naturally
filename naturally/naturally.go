// Implements "Naturally" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
	"strconv"
)

// Naturally implements sort.Interface by providing Less and
// using Len and Swap  methods of embedded []string
type Naturally struct {
	Val []string
}

// partition string into numeric and non-numeric parts
func partition(s string, ch chan<- string) {
	defer close(ch)
	numeric := false
	last := 0
	for ii, c := range s {
		if c >= '0' && c <= '9' {
			if numeric || last == ii {
				// either at start or already in numeric
				// value.  move on
				numeric = true
				continue
			}
			// end of non-numeric -- send back what we've got
			r := s[last:ii]
			ch <- r
			// numeric part starts at this char
			last = ii
			numeric = true
			continue
		}
		// non-numeric
		if !numeric || last == ii {
			numeric = false
			continue
		}
		// end of numeric
		r := s[last:ii]
		ch <- r
		// numeric part starts at next char
		last = ii
		numeric = false
		continue
	}
	ch <- s[last:]
}

func (p Naturally) Len() int {
	return len(p.Val)
}

func (p Naturally) Swap(a, b int) {
	p.Val[a], p.Val[b] = p.Val[b], p.Val[a]
}

func (p Naturally) Less(a, b int) bool {
	// part string -- numeric and non
	chA := make(chan string)
	chB := make(chan string)

	go partition(p.Val[a], chA)
	go partition(p.Val[b], chB)

	for {
		partA, okA := <-chA
		if !okA {
			// nothing more on A -- shorter or same as B
			return true
		}

		partB, okB := <-chB
		if !okB {
			// nothing more on B -- shorter than A
			return false
		}
		if partA == partB {
			// same -- move on
			continue
		}

		// not same string -- check numeric vals
		intA, errintA := strconv.Atoi(partA)
		intB, errintB := strconv.Atoi(partB)
		if errintA == nil {
			// A numeric
			if errintB == nil {
				// A & B numeric
				return intA < intB
			}
			return true
		}
		// A not numeric
		if errintB == nil {
			return false
		}
		// both string
		return partA < partB
	}
	// reached the end
	return true
}
