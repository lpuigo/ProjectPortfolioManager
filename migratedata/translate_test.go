package migratedata

import (
	"fmt"
	"os"
	"testing"
)

const (
	Dico     = "french2english.csv"
	PrjFile1 = "Projets Novagile formatted.json"
	PrjFile2 = "C:/Users/Laurent/Google Drive/Travail/NOVAGILE/Gouvernance/Ptf Projets/NovagileProjectManager/Ressources/Projets Novagile.xlsx.json"
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

	targertFile := PrjFile2 + ".translated"

	sf, err := os.Open(PrjFile2)
	if err != nil {
		t.Fatal("could not open source", err.Error())
	}
	defer sf.Close()

	tf, err := os.Create(targertFile)
	if err != nil {
		t.Fatal("could not open target", err.Error())
	}
	defer tf.Close()

	err = d.Translate(sf, tf)
	if err != nil {
		t.Fatal("TranslateFile returns", err.Error())
	}
}
