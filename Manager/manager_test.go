package Manager

import (
	"testing"
)

const (
	prjfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\Projets Novagile.xlsx.json`
	//StatFile = `C:\Users\Laurent\Google Drive\Golang\src\github.com\lpuig\Novagile\Ressources\Test Stats Projets Novagile.json`
	csvfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-03.csv`
)

func TestNewManager(t *testing.T) {
	_, err := NewManager(prjfile, StatFile)
	if err != nil {
		t.Fatal("NewManager returned", err.Error())
	}
}
