package RecordSet

import (
	"encoding/csv"
	"fmt"
	"io"
)

type RecordSet struct {
	header *Header
	data   []Record
}

// NewRecordSet returns an new empty RecordSet (both Header and Data are empty)
func NewRecordSet() *RecordSet {
	c := &RecordSet{data: []Record{}}
	return c
}

// CreateSubRecordSet returns an new empty (no data) RecordSet using same Header
func (c *RecordSet) CreateSubRecordSet() *RecordSet {
	r := NewRecordSet()
	r.header = c.header
	return r
}

func (c *RecordSet) AddHeader(record Record) error {
	c.header = NewHeader(record)
	return nil
}

func (c *RecordSet) GetHeader() *Header {
	return c.header
}

func (c *RecordSet) Add(record Record) int {
	num := len(c.data)
	c.data = append(c.data, record)
	return num
}

func (c *RecordSet) Get(num int) Record {
	return c.data[num]
}

func (c *RecordSet) Len() int {
	return len(c.data)
}

func (c *RecordSet) GetRecords() []Record {
	return c.data
}

// AddCSVDataFrom populates the RecordSet with Data from given reader (CSV formated data) (Header is created if missing)
func (c *RecordSet) AddCSVDataFrom(r io.Reader) error {
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
			if c.header == nil {
				c.AddHeader(record)
			} else if len(record) != len(c.header.keys) {
				return fmt.Errorf("Uncompatible RecordSet structure : %s vs %s", record, c.header.keys)
			} // else do nothing as reusing same header
		} else {
			c.data = append(c.data, record)
		}
		numline++
	}
	return nil
}

// WriteCSVTo writes all RecordSet records to given writer using CSVFormat (; delimitor, CRLF record separator)
func (c *RecordSet) WriteCSVTo(w io.Writer) error {
	csvw := csv.NewWriter(w)
	csvw.UseCRLF = true
	csvw.Comma = ';'

	err := csvw.Write(c.GetHeader().GetKeys())
	if err != nil {
		return err
	}

	for _, record := range c.data {
		err := csvw.Write(record)
		if err != nil {
			return err
		}
	}
	csvw.Flush()
	return nil
}
