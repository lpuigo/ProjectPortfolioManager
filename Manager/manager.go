package Manager

import (
	"encoding/json"
	"errors"
	"fmt"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/lpuig/Novagile/Model"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Manager struct {
	Projects *PrjManager
	Stats    *StatManager
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

func (m *Manager) GetPrjPtf(w io.Writer) {
	prjs := []*fm.Project{}
	m.Projects.RLock() // Ensure PTf is not modified while being cloned
	m.Stats.RLock()
	// TODO figure out how to avoid allocating full ptf clone before marshaling
	for _, p := range m.Projects.GetProjectsPtf().Projects {
		prjs = append(prjs, fm.CloneBEProject(p, m.Stats.HasStatsForProject(getProjectKey(p))))
	}
	m.Stats.RUnlock()
	m.Projects.RUnlock()
	json.NewEncoder(w).Encode(prjs)
}

func (m *Manager) GetPrjById(id int) *Model.Project {
	return m.Projects.GetProjectsPtf().GetPrjById(id)
}

func (m *Manager) UpdateProjectsSpentTime() {
	m.Projects.WLock()
	defer m.Projects.WUnlockWithPersist()
	m.Stats.RLock()
	defer m.Stats.RUnlock()

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

func (m *Manager) UpdateProject(op, np *Model.Project) bool {
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

func (m *Manager) CreateProject(p *Model.Project) (*Model.Project, bool) {
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

func (m *Manager) ReinitStatsFromDir(dir string) error {
	m.Stats.WLock()
	m.Stats.ClearStats()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Unable to browse Dir : %s", err.Error())
	}
	nbRecord := 0
	for _, file := range files {
		f, err := os.Open(dir + file.Name())
		if err != nil {
			return err
		}
		t0 := time.Now()
		added, err := m.UpdateStat(f)
		dur := time.Since(t0)
		if err != nil {
			return fmt.Errorf("UpdateStat error %s", err.Error())
		}
		log.Printf("Stats updated from '%s': %d record(s) added (took %v)\n", file.Name(), added, dur)
		nbRecord += added
	}
	m.Stats.WUnlockWithPersist()
	m.UpdateProjectsSpentTime()
	return nil
}

func (m *Manager) UpdateStat(r io.Reader) (int, error) {
	return m.Stats.UpdateFrom(r)
}

func getProjectKey(p *Model.Project) (string, string) {
	return p.Client, p.Name
}

func (m *Manager) GetProjectStatById(id int, w io.Writer) error {
	prj := m.Projects.GetProjectsPtf().GetPrjById(id)
	if prj == nil {
		return fmt.Errorf("Project id %d not found", id)
	}

	//Retrieve Project Stat :
	ps := fm.ProjectStat{}
	m.Stats.RLock()
	m.Projects.RLock()
	var err error
	ps.Issues, ps.Dates, ps.TimeSpent, ps.TimeRemaining, ps.TimeEstimated, err = m.Stats.GetProjectStatInfo(getProjectKey(prj))
	m.Projects.RUnlock()
	m.Stats.RUnlock()
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(ps)
	return nil
}

func (m *Manager) GetProjectStatProjectList(w io.Writer) error {
	//Retrieve Project Stat :
	m.Stats.RLock()
	prjlist := m.Stats.GetProjectStatList()
	m.Stats.RUnlock()
	return json.NewEncoder(w).Encode(fm.NewProjectStatNameFromList(prjlist, "!"))
}
