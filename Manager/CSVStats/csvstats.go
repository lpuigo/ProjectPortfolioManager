package CSVStats

import (
	"encoding/csv"
	"fmt"
	pcsv "github.com/lpuig/Novagile/Manager/ProcessCSV"
	"io"
	"os"
)

type IndexDesc struct {
	name string
	cols []string
}

func NewIndexDesc(name string, cols ...string) IndexDesc {
	return IndexDesc{name: name, cols: cols}
}

type CSVStats struct {
	data       *pcsv.CSVStats
	indexDescs []IndexDesc
	indexes    map[string]*index
}

func NewCSVStats(indexes ...IndexDesc) *CSVStats {
	c := &CSVStats{}
	c.data = pcsv.NewCSVStats()
	c.indexDescs = indexes
	c.indexes = map[string]*index{}
	return c
}

func (c *CSVStats) AddHeader(record []string) error {
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

func (c *CSVStats) Len() int {
	return c.data.Len()
}

// Add adds the given record, updating indexes
func (c *CSVStats) Add(record []string) {
	num := c.data.Add(record)
	for _, v := range c.indexes {
		v.Add(record, num)
	}
}

// GetKeys returns all registered keys for given idxname
func (c *CSVStats) GetKeys(idxname string) []string {
	if i, found := c.indexes[idxname]; found {
		return i.index.Keys()
	}
	return nil
}

// GetRecordKey returns the idxname related key for given record
func (c *CSVStats) GetRecordKey(idxname string, record []string) string {
	if i, found := c.indexes[idxname]; found {
		return i.genKey(record)
	}
	return ""
}

func (c *CSVStats) GetRecords(idxname, key string) [][]string {
	i, found := c.indexes[idxname]
	if !found {
		return nil
	}
	r, found := i.index[key]
	if !found {
		return nil
	}
	res := [][]string{}
	for _, pos := range r {
		res = append(res, c.data.Get(pos))
	}
	return res
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
			c.Add(record)
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
	subset := c.GetRecords(idxname, key)
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
