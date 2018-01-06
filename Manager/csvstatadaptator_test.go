package Manager

import (
	"testing"
	"time"
)

const (
	prjfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\Projets Novagile.xlsx.json`
	//StatFile = `C:\Users\Laurent\Google Drive\Golang\src\github.com\lpuig\Novagile\Ressources\Test Stats Projets Novagile.json`
	csvfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export_jira_modif.csv`
)

func TestUpdateStatPortfolioFromCSVFile(t *testing.T) {
	m, err := NewManager(prjfile, StatFile)
	if err != nil {
		t.Fatal("NewManager returned", err.Error())
	}
	err = m.UpdateStatFromCSVFile(csvfile)
	if err != nil {
		t.Error("UpdateStatFromCSVFile returned", err.Error())
	}
	time.Sleep(4 * time.Second)
}
