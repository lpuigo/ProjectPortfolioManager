package RecordIndexedSet

import (
	rs "github.com/lpuig/Novagile/Manager/RecordSet"
	"sort"
)

type Association struct {
	data     map[string]string
	genKey   rs.KeyGenerator
	genValue rs.KeyGenerator
}

func NewAssociation(keyfunc, valuefunc rs.KeyGenerator) *Association {
	a := &Association{}
	a.data = make(map[string]string)
	a.genKey, a.genValue = keyfunc, valuefunc
	return a
}

// Clone returns a new empty Association having the same Key / Value Generators (but no data)
func (a *Association) Clone() *Association {
	return NewAssociation(a.genKey, a.genValue)
}

// Set creates or replaces key entry with value
func (a *Association) Set(key, value string) {
	a.data[key] = value
}

// Get returns value associated with key, if exists, or def otherwise
func (a *Association) Get(key, def string) string {
	if v, found := a.data[key]; found {
		return v
	}
	return def
}

// Keys returns all known keys sorted
func (a *Association) Keys() []string {
	res := make([]string, len(a.data))
	i := 0
	for k, _ := range a.data {
		res[i] = k
		i++
	}
	sort.Strings(res)
	return res
}

func (a *Association) KeysMatching(value string) []string {
	res := []string{}
	for k, v := range a.data {
		if v == value {
			res = append(res, k)
		}
	}
	sort.Strings(res)
	return res
}

func (a *Association) Values() []string {
	vals := NewAssociation(nil, nil)
	for _, v := range a.data {
		vals.Set(v, "")
	}
	return vals.Keys()
}

func (a *Association) UpdateWith(r rs.Record) {
	a.Set(a.genKey(r), a.genValue(r))
}
