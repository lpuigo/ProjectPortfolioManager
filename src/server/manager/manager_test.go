package manager

import (
	"bytes"
	"fmt"
	"github.com/lpuig/prjptf/src/server/manager/fileprocesser"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

const (
	prjfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\novagile\Ressources\Projets Novagile.xlsx.json`
	//StatFile = `C:\Users\Laurent\Google Drive\Golang\src\github.com\lpuig\prjptf\Ressources\Test Stats Projets Novagile.json`
	csvfile = `C:\Users\Laurent\Golang\src\github.com\lpuig\novagile\Ressources\export Jira\extract 2018-01-03.csv`

	PrdStatFile        = `C:\Users\Laurent\Golang\src\github.com\lpuig\novagile\Ressources\Stats Projets Novagile.csv`
	UpdateStatDir      = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE\`
	ArchivedStatDir    = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Archived SRE`
	UpdateStatFile     = UpdateStatDir + `extract 2018-01-04.csv`
	BenchPrjFile       = `C:\Users\Laurent\Golang\src\github.com\lpuig\novagile\Ressources\Test\Bench\Projets Novagile.xlsx.json`
	BenchStatFile      = `C:\Users\Laurent\Golang\src\github.com\lpuig\novagile\Ressources\Test\Bench\Stats Projets Novagile.csv`
	BenchArchiveSREDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\novagile\Ressources\Test\Bench\ArchivedSRE`
	BenchExtractSREDir = `C:\Users\Laurent\Golang\src\github.com\lpuig\novagile\Ressources\Test\Bench\ExtractSRE`
	TestDBName         = ``
	TestDBUsrPwd       = ``
)

func TestNewManager(t *testing.T) {
	_, err := NewManager(prjfile, StatFile, TestDBUsrPwd, TestDBName)
	if err != nil {
		t.Fatal("NewManager returned", err.Error())
	}
}

func TestInitActualDataOnProdFile(t *testing.T) {
	m, err := NewManager(prjfile, PrdStatFile, TestDBUsrPwd, TestDBName)
	if err != nil {
		t.Fatalf("NewManager returns %s", err.Error())
	}
	m.Fp, err = fileprocesser.NewFileProcesser(UpdateStatDir, ArchivedStatDir)
	if err != nil {
		t.Fatalf("NewFileProcesser returns %s", err.Error())
	}

	w := new(bytes.Buffer)
	m.ReinitStats(w)

	fmt.Println(w.String())

	time.Sleep(4 * time.Second)
}

func TestManager_GetProjectStatById(t *testing.T) {
	m, err := NewManager(prjfile, PrdStatFile, TestDBUsrPwd, TestDBName)
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

	w := new(bytes.Buffer)
	m.GetProjectStatById(5, w)
	//for _, id := range prjIds {
	//	w := new(bytes.Buffer)
	//	m.GetProjectStatById(id, w)
	//}
}

func BenchmarkManager_GetProjectStatById(b *testing.B) {
	m, err := NewManager(prjfile, PrdStatFile, TestDBUsrPwd, TestDBName)
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

func BenchmarkManager_ReinitStats(b *testing.B) {
	m, err := NewManager(BenchPrjFile, BenchStatFile, TestDBUsrPwd, TestDBName)
	if err != nil {
		b.Fatalf("NewManager returns %s", err.Error())
	}
	m.Fp, err = fileprocesser.NewFileProcesser(BenchExtractSREDir, BenchArchiveSREDir)
	if err != nil {
		b.Fatalf("NewFileProcesser returns %s", err.Error())
	}

	f, err := os.Create("ReinitStats.pprof")
	if err != nil {
		b.Fatalf("could not create: %s", err.Error())
	}
	defer f.Close()
	pprof.StartCPUProfile(f)

	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			w := new(bytes.Buffer)
			m.ReinitStats(w)
			//b.Log(w.String())

		}
	}
	pprof.StopCPUProfile()

	time.Sleep(4 * time.Second)
}
