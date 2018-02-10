package RecordIndexedSet

import (
	rs "github.com/lpuig/Novagile/Manager/RecordSet"
	"sort"
)

type Association struct {
	values   map[string]string
	keys     map[string][]string
	genKey   rs.KeyGenerator
	genValue rs.KeyGenerator
}

func NewAssociation(keyfunc, valuefunc rs.KeyGenerator) *Association {
	a := &Association{}
	a.values = make(map[string]string)
	a.keys = make(map[string][]string)
	a.genKey, a.genValue = keyfunc, valuefunc
	return a
}

// Clone returns a new empty Association having the same Key / Value Generators (but no values)
func (a *Association) Clone() *Association {
	return NewAssociation(a.genKey, a.genValue)
}

// Set creates or replaces key entry with value
func (a *Association) Set(key, value string) {
	oldv, kfound := a.values[key]
	if kfound {
		if oldv == value { // nothing new here, exit
			return
		} else { // new value provided for key
			// remove key from keys[oldv]
			ks := a.keys[oldv]
			for i, k := range ks {
				if k == key {
					ks = append(ks[:i], ks[i+1:]...)
					break
				}
			}
			if len(ks) == 0 {
				delete(a.keys, oldv)
			} else {
				a.keys[oldv] = ks
			}
		}
	}
	// set new value
	a.values[key] = value
	listKeys, found := a.keys[value]
	if !found {
		a.keys[value] = []string{key}
	} else {
		a.keys[value] = append(listKeys, key)
	}
}

// Get returns value associated with key, if exists, or def otherwise
func (a *Association) Get(key, def string) string {
	if v, found := a.values[key]; found {
		return v
	}
	return def
}

// Keys returns all known keys sorted
func (a *Association) Keys() []string {
	res := make([]string, len(a.values))
	i := 0
	for k, _ := range a.values {
		res[i] = k
		i++
	}
	sort.Strings(res)
	return res
}

func (a *Association) KeysMatching(value string) []string {
	res := a.keys[value]
	sort.Strings(res)
	return res
}

func (a *Association) HasValue(value string) bool {
	_, found := a.keys[value]
	return found
}

func (a *Association) Values() []string {
	res := make([]string, len(a.keys))
	i := 0
	for k, _ := range a.keys {
		res[i] = k
		i++
	}
	sort.Strings(res)
	return res
}

func (a *Association) UpdateWith(r rs.Record) {
	a.Set(a.genKey(r), a.genValue(r))
}
