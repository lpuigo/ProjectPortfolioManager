package RecordIndexedSet

import (
	"strings"
	"testing"
)

func newAssociation() *Association {
	a := NewAssociation(nil, nil)
	a.Set("A", "0")
	a.Set("B", "2")
	a.Set("A", "1")
	a.Set("C", "1")
	return a
}

func TestAssociation_Get(t *testing.T) {
	a := newAssociation()

	tvl := []struct {
		key    string
		def    string
		expres string
	}{
		{"A", "", "1"},
		{"B", "C", "2"},
		{"Z", "Z", "Z"},
	}

	for _, tv := range tvl {
		res := a.Get(tv.key, tv.def)
		if res != tv.expres {
			t.Errorf("Association.Get('%s', '%s') returns %s instead of '%s'", tv.key, tv.def, res, tv.expres)
		}
	}
}

func TestAssociation_KeysValues(t *testing.T) {
	a := newAssociation()

	reskeys := strings.Join(a.Keys(), " ")
	exptres := "A B C"
	if reskeys != exptres {
		t.Errorf("Association.Keys() returns [%s] instead of [%s]", reskeys, exptres)
	}

	resvalues := strings.Join(a.Values(), " ")
	exptres = "1 2"
	if resvalues != exptres {
		t.Errorf("Association.Values() returns [%s] instead of [%s]", resvalues, exptres)
	}
}

func TestAssociation_KeysMatching(t *testing.T) {
	a := newAssociation()

	tvl := []struct {
		value  string
		expres string
	}{
		{"0", ""},
		{"1", "A C"},
		{"2", "B"},
		{"3", ""},
	}

	for _, tv := range tvl {
		res := strings.Join(a.KeysMatching(tv.value), " ")
		if res != tv.expres {
			t.Errorf("Association.KeysMatching('%s') returns %s instead of '%s'", tv.value, res, tv.expres)
		}
	}
}

func TestAssociation_HasValue(t *testing.T) {
	a := newAssociation()

	tvl := []struct {
		value  string
		expres bool
	}{
		{"0", false},
		{"1", true},
		{"2", true},
		{"3", false},
	}

	for _, tv := range tvl {
		res := a.HasValue(tv.value)
		if res != tv.expres {
			t.Errorf("Association.KeysMatching('%s') returns %s instead of '%s'", tv.value, res, tv.expres)
		}
	}
}
