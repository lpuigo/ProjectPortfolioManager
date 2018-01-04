package CSVStats

import (
	pcsv "github.com/lpuig/Novagile/Manager/ProcessCSV"
)

type MultiIndex map[string][]int

func NewMultiIdIndex() MultiIndex {
	return MultiIndex{}
}

func (i MultiIndex) Add(key string, pos int) {
	if !i.Has(key) {
		i[key] = []int{}
	}
	i[key] = append(i[key], pos)
}

func (i MultiIndex) Has(key string) bool {
	_, found := i[key]
	return found
}

func (i MultiIndex) Get(key string) ([]int, bool) {
	e, found := i[key]
	return e, found
}

type index struct {
	genKey pcsv.RecordSelector
	index  MultiIndex
}

func newIndex(rs pcsv.RecordSelector) *index {
	return &index{
		genKey: rs,
		index:  NewMultiIdIndex(),
	}
}

func (i *index) Add(record []string, num int) {
	i.index.Add(i.genKey(record), num)
}
