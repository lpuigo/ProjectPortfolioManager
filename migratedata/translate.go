package migratedata

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	rs "github.com/lpuig/novagile/manager/recordset"
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

func (d Dictionnary) TranslateFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {             // internally, it advances token based on separator
		orig := scanner.Text()

		for _, sd := range d.GetRecords() {
			orig = strings.Replace(orig, sd[0], sd[1], -1)
		}

		fmt.Println(orig)
	}
	return nil
}