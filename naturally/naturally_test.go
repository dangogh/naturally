package naturally_test

import (
	"github.com/user/naturally"
	"sort"
	"testing"
)

var testdata = [...]string{"A1", "A0", "A21", "A11", "A111", "A2"}
var testexp = [...]string{"A0", "A1", "A2", "A11", "A21", "A111"}

func TestSorts(t *testing.T) {
	n := naturally.Naturally{testdata[0:]}
	sort.Sort(n)

	for ii, got := range testdata {
		exp := testexp[ii]
		if got != exp {
			t.Errorf("Got %v; expected %v\n", got, exp)
		}
	}
}
