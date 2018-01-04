package Manager

import "testing"

const StatFile = `C:\Users\Laurent\Google Drive\Golang\src\github.com\lpuig\Novagile\Ressources\Test Stats Projets Novagile.json`

func TestInitStatManagerPersistFile(t *testing.T) {
	err := InitStatManagerPersistFile(StatFile)
	if err != nil {
		t.Errorf("InitStatManagerPersistFile: %s", err.Error())
	}
}

func TestNewStatManagerFromPersistFile(t *testing.T) {
	_, err := NewStatManagerFromPersistFile(StatFile)
	if err != nil {
		t.Errorf("NewStatManagerFromPersistFile: %s", err.Error())
	}
}
