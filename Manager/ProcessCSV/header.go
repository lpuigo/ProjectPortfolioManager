package ProcessCSV

import (
	"fmt"
)

type Header struct {
	keys  []string
	index map[string]int
}

type KeyGenerator func([]string) string

func NewHeader(record []string) *Header {
	h := &Header{
		keys:  record,
		index: map[string]int{},
	}
	for i, s := range record {
		h.index[s] = i
	}
	return h
}

func (h *Header) GetKeys() []string {
	return h.keys
}

func (h *Header) getNums(colname ...string) ([]int, error) {
	res := []int{}
	for _, s := range colname {
		n, found := h.index[s]
		if !found {
			return nil, fmt.Errorf("Column name error : '%s' is not in %s", s, h.keys)
		}
		res = append(res, n)
	}
	return res, nil
}

func (h *Header) NewKeyGenerator(colname ...string) (KeyGenerator, error) {
	colnums, err := h.getNums(colname...)
	if err != nil {
		return nil, err
	}
	return func(record []string) string {
		if len(colnums) == 0 {
			return "!"
		}
		res := ""
		for _, n := range colnums {
			res += "!" + record[n]
		}
		return res
	}, nil
}
