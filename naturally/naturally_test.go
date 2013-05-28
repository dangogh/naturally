package naturally_test

import (
	"github.com/user/naturally"
	"sort"
	"testing"
)

var testdata = [...]string{"ZZZ", "22", "Z22", "B22", "A0022", "A1"}
var testexp = [...]string{"", "22", "Z22", "B22", "A0022", "A1", "ZZZ"}

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
