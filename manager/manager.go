package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	fm "github.com/lpuig/novagile/client/frontmodel"
	fpr "github.com/lpuig/novagile/manager/fileprocesser"
	"github.com/lpuig/novagile/model"
	"io"
	"os"
	"time"
)

type Manager struct {
	Projects *PrjManager
	Stats    *StatManager
	Fp       *fpr.FileProcesser
}

func NewManager(prjfile, statfile string) (*Manager, error) {
	m := &Manager{}
	pm, err := NewPrjManagerFromPersistFile(prjfile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve Project Portfolio Data: %s", err.Error()))
	}
	m.Projects = pm

	sm, err := NewStatManagerFromFile(statfile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve Stat Portfolio Data: %s", err.Error()))
	}
	m.Stats = sm

	m.UpdateProjectsSpentTime()
	return m, nil
}

func (m *Manager) AddStatFileDirs(input, archive string) error {
	fp, err := fpr.NewFileProcesser(input, archive)
	if err != nil {
		return err
	}
	m.Fp = fp
	return nil
}

func (m *Manager) GetPrjPtf(w io.Writer) {
	prjs := make([]*fm.Project, 0)
	m.Projects.RLock() // Ensure PTf is not modified while being cloned
	m.Stats.RLock()
	for _, p := range m.Projects.GetProjectsPtf().Projects {
		prjs = append(prjs, fm.CloneBEProject(p, m.Stats.HasStatsForProject(getProjectKey(p))))
	}
	m.Stats.RUnlock()
	m.Projects.RUnlock()
	json.NewEncoder(w).Encode(prjs)
}

func (m *Manager) GetPrjById(id int) *model.Project {
	m.Projects.RLock()
	defer m.Projects.RUnlock()
	return m.Projects.GetProjectsPtf().GetPrjById(id)
}

func (m *Manager) UpdateProjectsSpentTime() {
	m.Stats.RLock()
	defer m.Stats.RUnlock()
	m.Projects.WLock()
	defer m.Projects.WUnlockWithPersist()

	m.updateProjectsSpentTime()
}

func (m *Manager) updateProjectsSpentTime() {
	for _, p := range m.Projects.GetProjectsPtf().Projects {
		hasStat := m.Stats.HasStatsForProject(getProjectKey(p))
		if hasStat {
			nWL, err := m.Stats.GetProjectSpentWL(getProjectKey(p))
			if err != nil {
				panic(err.Error())
			}
			p.CurrentWL = nWL
		}
	}
}

func (m *Manager) UpdateProject(op, np *model.Project) bool {
	m.Projects.WLock()
	defer m.Projects.WUnlockWithPersist()
	m.Stats.RLock()
	defer m.Stats.RUnlock()
	op.Update(np)
	hasStat := m.Stats.HasStatsForProject(getProjectKey(op))
	if hasStat {
		nWL, err := m.Stats.GetProjectSpentWL(getProjectKey(op))
		if err != nil {
			panic(err.Error())
		}
		op.CurrentWL = nWL
	}
	return hasStat
}

func (m *Manager) CreateProject(p *model.Project) (*model.Project, bool) {
	m.Projects.WLock()
	defer m.Projects.WUnlockWithPersist()
	m.Projects.GetProjectsPtf().AddPrj(p)
	m.Stats.RLock()
	defer m.Stats.RUnlock()
	hasStat := m.Stats.HasStatsForProject(getProjectKey(p))
	if hasStat {
		nWL, err := m.Stats.GetProjectSpentWL(getProjectKey(p))
		if err != nil {
			panic(err.Error())
		}
		p.CurrentWL = nWL
	}
	return p, hasStat
}

func (m *Manager) DeleteProject(id int) bool {
	m.Projects.WLock()
	found := m.Projects.GetProjectsPtf().DeletePrj(id)
	if found {
		m.Projects.WUnlockWithPersist()
	} else {
		m.Projects.WUnlock()
	}
	return found
}

func (m *Manager) GetProjectsPtfXLS(w io.Writer) {
	m.Projects.RLock()
	defer m.Projects.RUnlock()
	WritePortfolioToXLS(m.Projects.GetProjectsPtf(), w)
}

func (m *Manager) UpdateWithNewStatFiles(w io.Writer) error {
	m.Stats.WLock()
	defer m.Stats.WUnlockWithPersist()

	err := m.updateWithNewStatFiles(w)
	if err != nil {
		return err
	}

	m.Projects.WLock()
	defer m.Projects.WUnlockWithPersist()
	m.updateProjectsSpentTime()

	return nil
}

func (m *Manager) updateWithNewStatFiles(w io.Writer) error {
	tt := time.Now()
	err := m.Fp.ProcessAndArchive(func(file string) error {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(w, "%s", err.Error())
			return err
		}
		defer f.Close()
		t0 := time.Now()
		added, err := m.UpdateStat(f)
		dur := time.Since(t0)
		if err != nil {
			fmt.Fprintf(w, "%s", err.Error())
			return err
		}
		fmt.Fprintf(w, "Stats updated from '%s': %d record(s) added (took %v)\n", file, added, dur)
		return nil
	})
	if err != nil {
		fmt.Fprintf(w, "Update with new stat files aborted: %s", err.Error())
		return err
	}
	fmt.Fprintf(w, "Update with new stat files completed (took %v)", time.Since(tt))
	return nil
}

func (m *Manager) ReinitStats(w io.Writer) error {
	err := m.Fp.RestoreArchives()
	if err != nil {
		return err
	}
	m.Stats.WLock()
	defer m.Stats.WUnlockWithPersist()

	m.Stats.ClearStats()

	err = m.updateWithNewStatFiles(w)
	if err != nil {
		return err
	}

	m.Projects.WLock()
	defer m.Projects.WUnlockWithPersist()
	m.updateProjectsSpentTime()

	return nil
}

func (m *Manager) UpdateStat(r io.Reader) (int, error) {
	return m.Stats.UpdateFrom(r)
}

func getProjectKey(p *model.Project) (string, string) {
	return p.Client, p.Name
}

func (m *Manager) GetProjectStatById(id int, w io.Writer) error {
	prj := m.Projects.GetProjectsPtf().GetPrjById(id)
	if prj == nil {
		return fmt.Errorf("project id %d not found", id)
	}
	dates := prj.Situation.GetSituationToDate().DateListJSFormat()
	//Retrieve Project Stat :
	ps := fm.ProjectStat{}
	m.Stats.RLock()
	m.Projects.RLock()
	var err error
	c, n := getProjectKey(prj)
	sd, ed := "", ""
	if len(dates) > 0 {
		sd, ed = dates[0], dates[len(dates)-1]
	}
	ps.Issues, ps.Summaries, ps.StartDate, ps.TimeSpent, ps.TimeRemaining, ps.TimeEstimated, err = m.Stats.GetProjectStatInfoOnPeriod(c, n, sd, ed)
	m.Projects.RUnlock()
	m.Stats.RUnlock()
	if err != nil {
		return err
	}
	//easyjson.MarshalToWriter(ps, w)
	json.NewEncoder(w).Encode(ps)
	return nil
}

func (m *Manager) GetProjectStatProjectList(id int, w io.Writer) error {
	m.Projects.RLock()
	defer m.Projects.RUnlock()
	m.Stats.RLock()
	defer m.Stats.RUnlock()

	var prjlist []string
	if id != -1 {
		prj := m.Projects.GetProjectsPtf().GetPrjById(id)
		if prj == nil {
			return fmt.Errorf("project id %d not found", id)
		}

		prjlist = m.Stats.GetProjectStatListSortedBySimilarity(prj.Client+"!"+prj.Name, m.Projects.GetProjectsPtf().GetPrjClientName("!"))
	} else {
		prjlist = m.Stats.GetProjectStatList(m.Projects.GetProjectsPtf().GetPrjClientName("!"))
	}

	return json.NewEncoder(w).Encode(fm.NewProjectStatNameFromList(prjlist, "!"))
}
