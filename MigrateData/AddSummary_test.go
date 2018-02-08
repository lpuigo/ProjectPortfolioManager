package MigrateData

import (
	ris "github.com/lpuig/Novagile/Manager/RecordIndexedSet"
	"io"
	"os"
	"path/filepath"
	"testing"
)

const (
	JiraStatDir     = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE\`
	ArchivedStatDir = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Archived SRE\`
	MigratedStatDir = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Migrated SRE\`
)

func TestMigrateSet(t *testing.T) {
	// get target file
	file := JiraStatDir + "extract 2018-02-05.csv"
	tf, err := os.Open(file)
	if err != nil {
		t.Fatalf("Open returns: %s", err.Error())
	}

	ts, err := NewTargetSetFrom(tf)
	if err != nil {
		t.Fatalf("NewTargetSetFrom returns: %s", err.Error())
	}

	a := ts.GetLink("Issue-Summary")
	if a == nil {
		t.Fatalf("Link 'Issue-Summary' not found")
	}
	//for _, k := range a.Keys() {
	//	t.Logf("'%s' : '%s'", k, a.Get(k, "none"))
	//}
	err = ts.WriteCSVToFile(file + ".test")
	if err != nil {
		t.Fatalf("WriteCSVToFile returns: %s", err.Error())
	}
	os.Remove(file + ".test")
}

func getTargetModel(modelFile string, t *testing.T) *ris.RecordLinkedIndexedSet {
	tf, err := os.Open(modelFile)
	if err != nil {
		t.Fatalf("could not open: %s", err.Error())
	}

	ts, err := NewTargetSetFrom(tf)
	if err != nil {
		t.Fatalf("could not open create TargetSet: %s", err.Error())
	}
	return ts
}

func TestExtract(t *testing.T) {
	ts := getTargetModel(MigratedStatDir+"extract 2018-02-05.csv", t)

	archFile := ArchivedStatDir + "extract 2018-01-03.csv.zip"

	process := func(r io.Reader, file string) error {
		ns, err := MigrateSet(r, ts)
		if err != nil {
			return err
		}
		err = ns.WriteCSVToFile(filepath.Join(MigratedStatDir, file))
		if err != nil {
			return err
		}
		return nil
	}

	err := ExtractAndProcess(archFile, process)
	if err != nil {
		t.Fatalf("could not Extract: %s", err.Error())
	}
}

func TestMigrateAll(t *testing.T) {
	ts := getTargetModel(MigratedStatDir+"extract 2018-02-05.csv", t)

	migrateFile := func(r io.Reader, file string) error {
		ns, err := MigrateSet(r, ts)
		if err != nil {
			return err
		}
		err = ns.WriteCSVToFile(filepath.Join(MigratedStatDir, file))
		if err != nil {
			return err
		}
		return nil
	}

	processFile := func(archFile string) error {
		err := ExtractAndProcess(archFile, migrateFile)
		if err != nil {
			t.Errorf("could not Extract: %s", err.Error())
		}
		t.Logf("Processed %s", archFile)
		return nil
	}

	err := ProcessDir(ArchivedStatDir, processFile)
	if err != nil {
		t.Fatalf("could not Extract: %s", err.Error())
	}
}
