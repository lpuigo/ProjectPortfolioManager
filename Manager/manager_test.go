package Manager

import (
	"testing"
	"time"
)

const (
	prjfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\Projets Novagile.xlsx.json`
	//StatFile = `C:\Users\Laurent\Google Drive\Golang\src\github.com\lpuig\Novagile\Ressources\Test Stats Projets Novagile.json`
	csvfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-03.csv`

	PrdStatFile    = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\Stats Projets Novagile.csv`
	UpdateStatDir  = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE\`
	UpdateStatFile = UpdateStatDir + `extract 2018-01-04.csv`
)

func TestNewManager(t *testing.T) {
	_, err := NewManager(prjfile, StatFile)
	if err != nil {
		t.Fatal("NewManager returned", err.Error())
	}
}

func TestInitActualDataOnProdFile(t *testing.T) {
	m, err := NewManager(prjfile, PrdStatFile)
	if err != nil {
		t.Fatalf("NewManager returns %s", err.Error())
	}

	m.ReinitStatsFromDir(UpdateStatDir)

	time.Sleep(4 * time.Second)
}
