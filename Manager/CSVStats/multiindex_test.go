package CSVStats

import "testing"

func TestNewMultiIdIndex(t *testing.T) {
	i := newMultiIdIndex()
	i.Add("a", 0)
	i.Add("a", 1)
	i.Add("b", 1)

	_, f := i.Get("c")
	if f {
		t.Errorf("MultiIdIndex.Get('unknown key') returns true instead of false")
	}

	a, _ := i.Get("a")
	if len(a) != 2 {
		t.Errorf("MultiIdIndex.Get('known key') returns unexpected result length %v", a)
	}
	if a[0] != 0 || a[1] != 1 {
		t.Errorf("MultiIdIndex.Get('known key') returns unexpected result %v", a)
	}
}
