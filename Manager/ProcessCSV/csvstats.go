package ProcessCSV

import (
	"encoding/csv"
	"io"
)

type CSVStats struct {
	header *Header
	data   [][]string
}

func NewCSVStats() *CSVStats {
	c := &CSVStats{data: [][]string{}}
	return c
}

func (c *CSVStats) AddHeader(record []string) error {
	c.header = NewHeader(record)
	return nil
}

func (c *CSVStats) GetHeader() *Header {
	return c.header
}

func (c *CSVStats) Add(record []string) int {
	num := len(c.data)
	c.data = append(c.data, record)
	return num
}

func (c *CSVStats) Get(num int) []string {
	return c.data[num]
}

func (c *CSVStats) Len() int {
	return len(c.data)
}

func (c *CSVStats) GetRecords() [][]string {
	return c.data
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
			c.AddHeader(record)
		} else {
			c.data = append(c.data, record)
		}
		numline++
	}
	return nil
}

// WriteCSVTo writes all CSVStats records to given writer using CSVFormat (; delimitor, CRLF record separator)
func (c *CSVStats) WriteCSVTo(w io.Writer) error {
	csvw := csv.NewWriter(w)
	csvw.UseCRLF = true
	csvw.Comma = ';'

	return csvw.WriteAll(c.data)
}
