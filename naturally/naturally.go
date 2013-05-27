// Implements "Naturally" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
		"strconv"
		//"fmt"
		)

// Naturally implements sort.Interface by providing Less and 
// using Len and Swap  methods of embedded []string
type Naturally struct {
	Val []string
}

// partition string into numeric and non-numeric parts
func partition(s string, ch chan<- string) {
	numeric := false
	last := 0
	for ii, c := range s {
		if c >= '0' && c <= '9' {
			if numeric || last == ii {
				// either at the beginning or
				// already numeric -- move on
				numeric = true
				continue
			}
			r := s[last:ii]
			ch <- r
			// numeric part starts at next char
			last = ii + 1
			numeric = true
			continue
		}
		// non-numeric
		if !numeric || last == ii {
			numeric = false
			continue
		}
		r := s[last:ii]
		ch <- r
		// numeric part starts at next char
		last = ii + 1
		numeric = true
		continue
	}
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
		partA, errA := <-chA
		partB, errB := <-chB
		if !errA {
			// nothing more on A -- shorter or same as B
			return true
		}
		if !errB {
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
		if errintA != errintB {
			// if A numeric, A less else B less
			return errintA == nil
		}
		if errintA == nil {
			// both numeric: compare numerically
			if intA == intB {
				// same value -- leading 0's
				return len(partA) > len(partB)
			}
			return intA < intB
		}
		// both string
		return partA < partB
	}
	return true
}
