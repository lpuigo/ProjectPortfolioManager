package fileprocesser

import (
	"fmt"
	"testing"
)

const (
	inputDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\`
	archiveDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\ProcessedJiraFiles`

	prdInputDir   = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE`
	prdArchiveDir = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Archived SRE`
)

func TestFileProcesser_Process(t *testing.T) {
	fm, err := NewFileProcesser(inputDir, archiveDir)
	if err != nil {
		t.Fatalf("NewFileProcesser", err.Error())
	}
	err = fm.ProcessAndArchive(func(f string) error {
		fmt.Printf("Processing  %s\n", f)
		return nil
	})
	if err != nil {
		t.Error("ProcessAndArchive", err.Error())
	}
	//time.Sleep(time.Second)
}

func TestFileProcesser_RestoreArchives(t *testing.T) {
	fm, err := NewFileProcesser(inputDir, archiveDir)
	if err != nil {
		t.Fatalf("NewFileProcesser", err.Error())
	}
	err = fm.RestoreArchives()
	if err != nil {
		t.Error("RestoreArchives", err.Error())
	}
}

func TestFileProcesser_RestoreArchivesPROD(t *testing.T) {
	fm, err := NewFileProcesser(prdInputDir, prdArchiveDir)
	if err != nil {
		t.Fatalf("NewFileProcesser", err.Error())
	}
	err = fm.RestoreArchives()
	if err != nil {
		t.Error("RestoreArchives", err.Error())
	}
}
