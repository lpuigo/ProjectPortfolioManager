package CSVStats

import (
	pcsv "github.com/lpuig/Novagile/Manager/ProcessCSV"
)

type multiIndex map[string][]int

func newMultiIdIndex() multiIndex {
	return multiIndex{}
}

func (i multiIndex) Add(key string, pos int) {
	if !i.Has(key) {
		i[key] = []int{}
	}
	i[key] = append(i[key], pos)
}

func (i multiIndex) Has(key string) bool {
	_, found := i[key]
	return found
}

func (i multiIndex) Keys() []string {
	res := []string{}
	for k, _ := range i {
		res = append(res, k)
	}
	return res
}

func (i multiIndex) Get(key string) ([]int, bool) {
	e, found := i[key]
	return e, found
}

type index struct {
	genKey pcsv.KeyGenerator
	index  multiIndex
}

func newIndex(rs pcsv.KeyGenerator) *index {
	return &index{
		genKey: rs,
		index:  newMultiIdIndex(),
	}
}

func (i *index) Add(record []string, num int) {
	i.index.Add(i.genKey(record), num)
}
