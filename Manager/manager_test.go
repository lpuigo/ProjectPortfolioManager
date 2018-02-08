package Manager

import (
	"bytes"
	"fmt"
	"github.com/lpuig/Novagile/Manager/FileProcesser"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

const (
	prjfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\Projets Novagile.xlsx.json`
	//StatFile = `C:\Users\Laurent\Google Drive\Golang\src\github.com\lpuig\Novagile\Ressources\Test Stats Projets Novagile.json`
	csvfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-03.csv`

	PrdStatFile     = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\Stats Projets Novagile.csv`
	UpdateStatDir   = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE\`
	ArchivedStatDir = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Archived SRE`
	UpdateStatFile  = UpdateStatDir + `extract 2018-01-04.csv`
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
	m.Fp, err = FileProcesser.NewFileProcesser(UpdateStatDir, ArchivedStatDir)
	if err != nil {
		t.Fatalf("NewFileProcesser returns %s", err.Error())
	}

	w := new(bytes.Buffer)
	m.ReinitStats(w)

	fmt.Println(w.String())

	time.Sleep(4 * time.Second)
}

func TestManager_GetProjectStatById(t *testing.T) {
	m, err := NewManager(prjfile, PrdStatFile)
	if err != nil {
		t.Fatalf("could not create new manager: %s", err.Error())
	}

	prjIds := map[string]int{}
	for _, p := range m.Projects.GetProjectsPtf().Projects {
		pk := p.Client + "!" + p.Name
		if m.Stats.HasStatsForProject(getProjectKey(p)) {
			prjIds[pk] = p.Id
		}
	}

	for _, id := range prjIds {
		w := new(bytes.Buffer)
		m.GetProjectStatById(id, w)
	}
}

func BenchmarkManager_GetProjectStatById(b *testing.B) {
	m, err := NewManager(prjfile, PrdStatFile)
	if err != nil {
		b.Fatalf("could not create new manager: %s", err.Error())
	}

	prjIds := map[string]int{}
	for _, p := range m.Projects.GetProjectsPtf().Projects {
		pk := p.Client + "!" + p.Name
		if m.Stats.HasStatsForProject(getProjectKey(p)) {
			prjIds[pk] = p.Id
		}
	}

	f, err := os.Create("GetProjectStatById.pprof")
	if err != nil {
		b.Fatalf("could not create: %s", err.Error())
	}
	defer f.Close()
	pprof.StartCPUProfile(f)

	time.Sleep(3 * time.Second)
	for p, id := range prjIds {
		b.Run(fmt.Sprintf("Stat on %s", p), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				w := new(bytes.Buffer)
				m.GetProjectStatById(id, w)
			}
		})
	}
	pprof.StopCPUProfile()

}
