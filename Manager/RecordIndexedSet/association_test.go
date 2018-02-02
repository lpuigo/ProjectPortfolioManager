package RecordIndexedSet

import (
	"testing"
	"strings"
)

func TestNewAssociation(t *testing.T) {
	a := NewAssociation(nil, nil)
	a.Set("A", "0")
	a.Set("B", "2")
	a.Set("A", "1")


	tvl := []struct {
		key string
		def string
		expres string
	}{
		{"A", "", "1"},
		{"B", "C", "2"},
		{"C", "D", "D"},
	}

	for _, tv := range tvl {
		res := a.Get(tv.key, tv.def)
		if res != tv.expres {
			t.Errorf("Association.Get('%s', '%s') returns %s instead of '%s'", tv.key, tv.def, res,tv.expres)
		}
	}

	reskeys := strings.Join(a.Keys(), " ")
	if reskeys != "A B" {
		t.Errorf("Association.Keys() returns [%s] instead of [%s]", reskeys, "A B")
	}

	resvalues := strings.Join(a.Values(), " ")
	if resvalues != "1 2" {
		t.Errorf("Association.Values() returns [%s] instead of [%s]", resvalues, "1 2")
	}
}