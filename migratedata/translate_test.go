package migratedata

import (
	"fmt"
	"testing"
)

const (
	Dico = "french2english.csv"
	PrjFile1 = "Projets Novagile formatted.json"
	PrjFile2 = "Projets Novagile.xlsx.json"
)

func NewDico(t *testing.T) *Dictionnary {
	d, err := NewDictionnaryFromCSVFile(Dico)
	if err != nil {
		t.Fatal("NewDictionnaryFromCSVFile returns:", err.Error())
	}
	return d
}

func TestNewDictionnaryFromCSVFile(t *testing.T) {
	d := NewDico(t)

	for _, sd := range d.GetRecords() {
		fmt.Println(sd)
	}
}

func TestDictionnary_TranslateFile(t *testing.T) {
	d := NewDico(t)

	err :=d.TranslateFile(PrjFile2)
	if err != nil {
		t.Fatal("TranslateFile returns", err.Error())
	}
}
