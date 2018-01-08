package CSVStats

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

type CSVStats struct {
	data       *rs.RecordSet
	indexDescs []IndexDesc
	indexes    map[string]*index
}

func NewCSVStats(indexes ...IndexDesc) *CSVStats {
	c := &CSVStats{}
	c.data = rs.NewRecordSet()
	c.indexDescs = indexes
	c.indexes = map[string]*index{}
	return c
}

func (c *CSVStats) Len() int {
	return c.data.Len()
}

func (c *CSVStats) AddHeader(record rs.Record) error {
	c.data.AddHeader(record)
	h := c.data.GetHeader()
	for _, id := range c.indexDescs {
		rs, err := h.NewKeyGenerator(id.cols...)
		if err != nil {
			return fmt.Errorf("CSVStats Unable to create index '%s' : %s", id.name, err.Error())
		}
		c.indexes[id.name] = newIndex(rs)
	}
	return nil
}

// AddRecord adds the given record, updating indexes
func (c *CSVStats) AddRecord(record rs.Record) {
	num := c.data.Add(record)
	for _, v := range c.indexes {
		v.Add(record, num)
	}
}

// HasIndexKey returns true if idxname index has given key, false otherwise
func (c *CSVStats) HasIndexKey(idxname, key string) bool {
	i, found := c.indexes[idxname]
	if !found {
		return false
	}
	_, found = i.index[key]
	return found
}

// GetIndexesNames returns names of active indexes
func (c *CSVStats) GetIndexesNames() rs.Record {
	res := rs.Record{}
	for n, _ := range c.indexes {
		res = append(res, n)
	}
	sort.Strings(res)
	return res
}

// GetIndexKeys returns all registered keys for given idxname
func (c *CSVStats) GetIndexKeys(idxname string) rs.Record {
	if i, found := c.indexes[idxname]; found {
		return i.index.Keys()
	}
	return nil
}

func (c *CSVStats) GetKeyGeneratorByIndexDesc(compare IndexDesc) (rs.KeyGenerator, error) {
	return c.data.GetHeader().NewKeyGenerator(compare.cols...)
}

// GetRecordKeyByIndex returns key related to given record key using named index
func (c *CSVStats) GetRecordKeyByIndex(idxname string, record rs.Record) string {
	if i, found := c.indexes[idxname]; found {
		return i.genKey(record)
	}
	return ""
}

func (c *CSVStats) GetRecordsByIndexKey(idxname, key string) []rs.Record {
	i, found := c.indexes[idxname]
	if !found {
		return nil
	}
	r, found := i.index[key]
	if !found {
		return nil
	}
	res := []rs.Record{}
	for _, pos := range r {
		res = append(res, c.data.Get(pos))
	}
	return res
}

func (c *CSVStats) GetRecords() []rs.Record {
	return c.data.GetRecords()
}

// AddCSVDataFrom populates the CSVStats with Data from given reader (CSV formated data) (Header and Indexes are populated)
func (c *CSVStats) AddCSVDataFrom(r io.Reader) error {
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

func (c *CSVStats) AddCSVDataFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.AddCSVDataFrom(f)
}

// Max returns the bigger record from (idxname, key) subset according to compare's Cols key(s)
func (c *CSVStats) Max(idxname, key string, compare IndexDesc) []string {
	comp, err := c.data.GetHeader().NewKeyGenerator(compare.cols...)
	if err != nil {
		panic(fmt.Sprintf("CSVStats.Max: %s", err.Error()))
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

// WriteCSVTo writes all CSVStats records to given writer using CSVFormat (; delimitor, CRLF record separator)
func (c *CSVStats) WriteCSVTo(w io.Writer) error {
	return c.data.WriteCSVTo(w)
}

func (c *CSVStats) WriteCSVToFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.WriteCSVTo(f)
}
