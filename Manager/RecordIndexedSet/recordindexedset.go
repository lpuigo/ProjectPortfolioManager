package RecordIndexedSet

import (
	"encoding/csv"
	"fmt"
	rs "github.com/lpuig/Novagile/Manager/RecordSet"
	"io"
	"os"
	"sort"
)

type IndexDesc struct {
	name string
	cols []string
}

func NewIndexDesc(name string, cols ...string) IndexDesc {
	return IndexDesc{name: name, cols: cols}
}

type RecordIndexedSet struct {
	data       *rs.RecordSet
	indexDescs []IndexDesc
	indexes    map[string]*index
}

func NewRecordIndexedSet(indexes ...IndexDesc) *RecordIndexedSet {
	c := &RecordIndexedSet{}
	c.data = rs.NewRecordSet()
	c.indexDescs = indexes
	c.indexes = map[string]*index{}
	return c
}

// CreateSubSet returns an new empty (no data) RecordIndexedSet using same Header
func (c *RecordIndexedSet) CreateSubSet(indexes ...IndexDesc) (*RecordIndexedSet, error) {
	r := NewRecordIndexedSet(indexes...)
	err := r.AddHeader(c.data.GetHeader().GetKeys())
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *RecordIndexedSet) Len() int {
	return c.data.Len()
}

func (c *RecordIndexedSet) GetRecordColNumByName(colname ...string) ([]int, error) {
	return c.data.GetRecordColNumByName(colname...)
}

func (c *RecordIndexedSet) AddHeader(record rs.Record) error {
	c.data.AddHeader(record)
	h := c.data.GetHeader()
	for _, id := range c.indexDescs {
		rs, err := h.NewKeyGenerator(id.cols...)
		if err != nil {
			return fmt.Errorf("RecordIndexedSet Unable to create index '%s' : %s", id.name, err.Error())
		}
		c.indexes[id.name] = newIndex(rs)
	}
	return nil
}

// AddRecord adds the given records, updating indexes
func (c *RecordIndexedSet) AddRecord(record ...rs.Record) {
	nums := c.data.Add(record...)
	for _, v := range c.indexes {
		v.Add(record, nums)
	}
}

// AddRecord adds the given records, updating indexes
func (c *RecordIndexedSet) AddRecords(records []rs.Record) {
	c.AddRecord(records...)
}

// HasIndexKey returns true if idxname index has given key, false otherwise
func (c *RecordIndexedSet) HasIndexKey(idxname, key string) bool {
	i, found := c.indexes[idxname]
	if !found {
		return false
	}
	_, found = i.index[key]
	return found
}

// GetIndexesNames returns names of active indexes
func (c *RecordIndexedSet) GetIndexesNames() rs.Record {
	res := rs.Record{}
	for n, _ := range c.indexes {
		res = append(res, n)
	}
	sort.Strings(res)
	return res
}

// GetIndexKeys returns all registered keys for given idxname
func (c *RecordIndexedSet) GetIndexKeys(idxname string) rs.Record {
	if i, found := c.indexes[idxname]; found {
		return i.index.Keys()
	}
	return nil
}

// GetIndexKeysByPrefix returns all registered keys for given idxname having given prefix
func (c *RecordIndexedSet) GetIndexKeysByPrefix(idxname, prefix string) rs.Record {
	if i, found := c.indexes[idxname]; found {
		return i.index.KeysByPrefix(prefix)
	}
	return nil
}

func (c *RecordIndexedSet) GetKeyGeneratorByIndexDesc(compare IndexDesc) (rs.KeyGenerator, error) {
	return c.data.GetHeader().NewKeyGenerator(compare.cols...)
}

// GetRecordKeyByIndex returns key related to given record key using named index
func (c *RecordIndexedSet) GetRecordKeyByIndex(idxname string, record rs.Record) string {
	if i, found := c.indexes[idxname]; found {
		return i.genKey(record)
	}
	return ""
}

// GetRecordsByIndexKey returns all records related to given key on given idxname (nil if idxname or key not found)
func (c *RecordIndexedSet) GetRecordsByIndexKey(idxname, key string) []rs.Record {
	r := c.GetRecordNumsByIndexKey(idxname, key)
	if r == nil {
		return nil
	}
	res := []rs.Record{}
	for _, pos := range r {
		res = append(res, c.data.Get(pos))
	}
	return res
}

// GetRecordsByIndexKey returns all records related to given key on given idxname (nil if idxname or key not found)
func (c *RecordIndexedSet) GetRecordNumsByIndexKey(idxname, key string) []int {
	i, found := c.indexes[idxname]
	if !found {
		return nil
	}
	r, found := i.index[key]
	if !found {
		return nil
	}
	return r
}

// GetRecords returns slice of all RIS records
func (c *RecordIndexedSet) GetRecords() []rs.Record {
	return c.data.GetRecords()
}

// GetRecordByNum returns RIS record with given num
func (c *RecordIndexedSet) GetRecordByNum(num int) rs.Record {
	return c.data.Get(num)
}

// GetRecordByNum returns RIS record with given num
func (c *RecordIndexedSet) GetRecordsByNums(nums []int) []rs.Record {
	res := make([]rs.Record, len(nums))
	for i, n := range nums {
		res[i] = c.data.Get(n)
	}
	return res
}

// AddCSVDataFrom populates the RecordIndexedSet with Data from given reader (CSV formated data) (Header and Indexes are populated)
func (c *RecordIndexedSet) AddCSVDataFrom(r io.Reader) error {
	csvr := csv.NewReader(r)
	csvr.Comma = ';'
	csvr.Comment = '#'

	var numline int = 0
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if numline == 0 {
			err = c.AddHeader(record)
			if err != nil {
				return err
			}
		} else {
			c.AddRecord(record)
		}

		numline++
	}
	return nil
}

func (c *RecordIndexedSet) AddCSVDataFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.AddCSVDataFrom(f)
}

// Max returns the bigger record from (idxname, key) subset according to compare's Cols key(s)
func (c *RecordIndexedSet) Max(idxname, key string, compare IndexDesc) []string {
	comp, err := c.data.GetHeader().NewKeyGenerator(compare.cols...)
	if err != nil {
		panic(fmt.Sprintf("RecordIndexedSet.Max: %s", err.Error()))
	}
	subset := c.GetRecordsByIndexKey(idxname, key)
	if len(subset) == 0 {
		return nil
	}
	if len(subset) == 1 {
		return subset[0]
	}
	max := 0
	for j := 1; j < len(subset); j++ {
		if comp(subset[max]) < comp(subset[j]) {
			max = j
		}
	}
	return subset[max]
}

// WriteCSVTo writes all RecordIndexedSet records to given writer using CSVFormat (; delimitor, CRLF record separator)
func (c *RecordIndexedSet) WriteCSVTo(w io.Writer) error {
	return c.data.WriteCSVTo(w)
}

func (c *RecordIndexedSet) WriteCSVToFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.WriteCSVTo(f)
}
