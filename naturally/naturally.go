// Implements "Naturally" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
	"strconv"
	"fmt"
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
		fmt.Printf("===== Got %c at %v\n", c, ii)
		if c >= '0' && c <= '9' {
			fmt.Printf("======= is numeric\n" )
			if numeric || last == ii {
				// either at start or already in numeric
				// value.  Move on already numeric -- move on
				numeric = true
				continue
			}
			// end of numeric -- send back what we've got
			r := s[last:ii]
			fmt.Printf("===== Got a %v\n", r)
			ch <- r
			// numeric part starts at next char
			last = ii + 1
			numeric = false
			continue
		}
		// non-numeric
		if !numeric || last == ii {
			numeric = false
			continue
		}
		// end of non-numeric
		r := s[last:ii]
		ch <- r
		// numeric part starts at next char
		last = ii + 1
		numeric = true
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
	//defer close(chA)
	chB := make(chan string)
	//defer close(chB)

	fmt.Printf("=== Comparing %v (%v) with %v (%v)\n",
		p.Val[a], a, p.Val[b], b)
	go partition(p.Val[a], chA)
	go partition(p.Val[b], chB)

	for {
		fmt.Println("Start of loop")
		partA, errA := <-chA
		fmt.Printf("===  partA: %v errA: %v\n", partA, errA)
		
		partB, errB := <-chB
		fmt.Printf("===  partB: %v errB: %v\n", partB, errB)
		if errA {
			// nothing more on A -- shorter or same as B
			return true
		}
		if errB {
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
