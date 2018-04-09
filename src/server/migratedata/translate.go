package migratedata

import (
	"bufio"
	"io"
	"os"
	"strings"

	rs "github.com/lpuig/novagile/src/server/manager/recordset"
)

type Dictionnary struct {
	*rs.RecordSet
}

func NewDictionnaryFromCSVFile(file string) (*Dictionnary, error) {
	d := &Dictionnary{rs.NewRecordSet()}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = d.AddCSVDataFrom(f)
	return d, err
}

func (d Dictionnary) Translate(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {             // internally, it advances token based on separator
		orig := scanner.Text()

		for _, sd := range d.GetRecords() {
			orig = strings.Replace(orig, sd[0], sd[1], -1)
		}

		_, err := w.Write([]byte(orig))
		if err != nil {
			return err
		}
	}
	return nil
}