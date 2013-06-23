package naturally_test

import (
	"github.com/dangogh/naturally"
	"sort"
	"testing"
	"fmt"
)

var testdata0 = [...]string{"A1", "A0", "A21", "A11", "A111", "A2"}
var testexp0 = [...]string{"A0", "A1", "A2", "A11", "A21", "A111"}
/*
var testdata1 = [...]string{"1", "2", "30", "22", "0", "00", "3"}
var testexp1 = [...]string{"0", "00", "1", "2", "3", "22", "30"}

var testdata2 = [...]string{"A1", "A0", "A21", "A11", "A111", "A2"}
var testexp2 = [...]string{"A0", "A1", "A2", "A11", "A21", "A111"}

var testdata3 = [...]string{"A2A", "A0", "A1", "A", "AA0", "A0A", "A11A"}
var testexp3 = [...]string{"A", "A0", "A0A", "A1", "A2A", "A11A", "AA0"}
*/
func _testSorts(t *testing.T, data, expected []string) {
	n := naturally.Naturally{data[0:]}
	sort.Sort(n)

	for ii, got := range data {
		exp := expected[ii]
		if got != exp {
			t.Errorf("Got %v; expected %v\n", got, exp)
			fmt.Printf("Got %v; expected %v\n", got, exp)
			
		}
	}
}

func TestSortsA(t *testing.T) {
	fmt.Println("TestA")
	_testSorts(t, testdata0[0:], testexp0[0:])
}
/*
func TestSortsB(t *testing.T) {
	fmt.Println("TestB")
	_testSorts(t, testdata1[0:], testexp1[0:])
}
func TestSortsC(t *testing.T) {
	fmt.Println("TestC")
	_testSorts(t, testdata2[0:], testexp2[0:])
}
func TestSortsD(t *testing.T) {
	fmt.Println("TestD")
	_testSorts(t, testdata3[0:], testexp3[0:])
}
*/
