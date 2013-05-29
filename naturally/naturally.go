// Implements "Naturally" sort -- alphabetic portion sorted
// alphabetically, numeric part sorted numerically.
package naturally

import (
	"fmt"
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
	//defer close(chA)
	chB := make(chan string)
	//defer close(chB)

	fmt.Printf("=== Comparing %v (%v) with %v (%v)\n",
		p.Val[a], a, p.Val[b], b)
	go partition(p.Val[a], chA)
	go partition(p.Val[b], chB)

	for {
		fmt.Println("Start of loop")
		partA, okA := <-chA
		fmt.Printf("===  partA: %v okA: %v\n", partA, okA)
		if !okA {
			// nothing more on A -- shorter or same as B
			fmt.Printf("===== nothing more on partA\n")
			return true
		}

		partB, okB := <-chB
		fmt.Printf("===  partB: %v okB: %v\n", partB, okB)
		if !okB {
			// nothing more on B -- shorter than A
			fmt.Printf("===== nothing more on partB\n")
			return false
		}
		if partA == partB {
			// same -- move on
			fmt.Printf("===== partA == partB\n")
			continue
		}

		// not same string -- check numeric vals
		intA, errintA := strconv.Atoi(partA)
		intB, errintB := strconv.Atoi(partB)
		if errintA != errintB {
			// if A numeric, A less else B less
			fmt.Printf("===== only one of partA, partB numeric\n")
			return errintA == nil
		}
		if errintA == nil {
			fmt.Printf("===== compare numerically\n")
			// both numeric: compare numerically
			if intA == intB {
				// same value -- leading 0's
				fmt.Printf("===== same value\n")
				return len(partA) > len(partB)
			}
			fmt.Printf("===== a < b? %v\n", (intA<intB))
			return intA < intB
		}
		// both string
		fmt.Printf("===== a before b? %v\n", (partA < partB))
		return partA < partB
	}
	return true
}
