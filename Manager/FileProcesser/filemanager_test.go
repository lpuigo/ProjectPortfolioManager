package FileProcesser

import (
	"fmt"
	"testing"
)

const (
	inputDir   = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\`
	archiveDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\ProcessedJiraFiles`
)

func TestFileProcesser_Process(t *testing.T) {
	fm, err := NewFileManager(inputDir, archiveDir)
	if err != nil {
		t.Fatalf("NewFileManager", err.Error())
	}
	err = fm.Process(func(f string) error {
		fmt.Printf("Processing  %s\n", f)
		return nil
	})
	if err != nil {
		t.Error("Process", err.Error())
	}

}
