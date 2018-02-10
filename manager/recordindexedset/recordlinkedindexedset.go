package recordindexedset

import (
	"encoding/csv"
	"fmt"
	rs "github.com/lpuig/novagile/manager/recordset"
	"io"
	"os"
)

type LinkDesc struct {
	name       string
	keyIndex   string
	valueIndex string
}

func NewLinkDesc(name, keyIndex, valueIndex string) LinkDesc {
	return LinkDesc{name: name, keyIndex: keyIndex, valueIndex: valueIndex}
}

type RecordLinkedIndexedSet struct {
	*RecordIndexedSet
	linkDescs []LinkDesc
	links     map[string]*Association
}

func NewRecordLinkedIndexedSet(indexes ...IndexDesc) *RecordLinkedIndexedSet {
	rlis := &RecordLinkedIndexedSet{RecordIndexedSet: NewRecordIndexedSet(indexes...)}
	rlis.linkDescs = []LinkDesc{}
	rlis.links = make(map[string]*Association)
	return rlis
}

// GetLink returns link with given name, or nil if not found
func (c *RecordLinkedIndexedSet) GetLink(linkname string) *Association {
	if a, found := c.links[linkname]; found {
		return a
	}
	return nil
}

// AddLink adds future link with given name on given key/value indexes. Fails if given index names are not already declared
func (c *RecordLinkedIndexedSet) AddLink(linkdesc LinkDesc) error {
	checkIndex := func(indexname string) bool {
		for _, id := range c.indexDescs {
			if indexname == id.name {
				return true
			}
		}
		return false
	}

	if !checkIndex(linkdesc.keyIndex) {
		return fmt.Errorf("index '%s' not found", linkdesc.keyIndex)
	}
	if !checkIndex(linkdesc.valueIndex) {
		return fmt.Errorf("index '%s' not found", linkdesc.valueIndex)
	}

	c.linkDescs = append(c.linkDescs, linkdesc)

	return nil
}

func (c *RecordLinkedIndexedSet) AddHeader(record rs.Record) error {
	err := c.RecordIndexedSet.AddHeader(record)
	if err != nil {
		return err
	}
	for _, ld := range c.linkDescs {
		kg := c.RecordIndexedSet.indexes[ld.keyIndex].genKey
		vg := c.RecordIndexedSet.indexes[ld.valueIndex].genKey
		c.links[ld.name] = NewAssociation(kg, vg)
	}
	return nil
}

// CreateSubSet returns an new empty (no values) RecordIndexedSet using same Header
func (c *RecordLinkedIndexedSet) CreateSubSet(indexes []IndexDesc, links []LinkDesc) (*RecordLinkedIndexedSet, error) {
	r := NewRecordLinkedIndexedSet(indexes...)
	for _, ld := range links {
		err := r.AddLink(ld)
		if err != nil {
			return nil, err
		}
	}
	err := r.AddHeader(c.data.GetHeader().GetKeys())
	if err != nil {
		return nil, err
	}
	return r, nil
}

// AddRecord adds the given records, updating indexes
func (c *RecordLinkedIndexedSet) AddRecord(record ...rs.Record) {
	c.RecordIndexedSet.AddRecord(record...)
	for _, r := range record {
		for _, l := range c.links {
			l.UpdateWith(r)
		}
	}
}

// AddRecord adds the given records, updating indexes
func (c *RecordLinkedIndexedSet) AddRecords(records []rs.Record) {
	c.AddRecord(records...)
}

// AddCSVDataFrom populates the RecordLinkedIndexedSet with Data from given reader (CSV formated values) (Header, Indexes and Links are populated)
func (c *RecordLinkedIndexedSet) AddCSVDataFrom(r io.Reader) error {
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

func (c *RecordLinkedIndexedSet) AddCSVDataFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.AddCSVDataFrom(f)
}
