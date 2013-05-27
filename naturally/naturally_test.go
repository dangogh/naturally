
package naturally_test

import (
	"testing"
	"github.com/user/naturally"
	"sort"
	)

var testdata = [...]string{"", "22", "Z22", "B22", "A0022", "A1"}
var testexp = [...]string{"", "22", "Z22", "B22", "A0022", "A1"}

func TestSorts(t *testing.T) {
	n := naturally.Naturally{testdata[0:]}
	sort.Sort( n )
}
