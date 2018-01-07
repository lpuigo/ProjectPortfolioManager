package Manager

import (
	"encoding/json"
	"errors"
	"fmt"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/lpuig/Novagile/Model"
	"io"
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
	m.Projects.RUnlock()
	m.Projects.RUnlock()
	json.NewEncoder(w).Encode(prjs)
}

func (m *Manager) GetPrjById(id int) *Model.Project {
	return m.Projects.GetProjectsPtf().GetPrjById(id)
}

func (m *Manager) UpdateProject(op, np *Model.Project) bool {
	m.Projects.WLock()
	defer m.Projects.WUnlock()
	m.Stats.RLock()
	defer m.Stats.RUnlock()
	op.Update(np)
	return m.Stats.HasStatsForProject(getProjectKey(np))
}

func (m *Manager) CreateProject(p *Model.Project) (*Model.Project, bool) {
	m.Projects.WLock()
	defer m.Projects.WUnlock()
	m.Projects.GetProjectsPtf().AddPrj(p)
	m.Stats.RLock()
	defer m.Stats.RUnlock()
	return p, m.Stats.HasStatsForProject(getProjectKey(p))
}

func (m *Manager) DeleteProject(id int) bool {
	m.Projects.WLock()
	found := m.Projects.GetProjectsPtf().DeletePrj(id)
	if found {
		m.Projects.WUnlock()
	} else {
		m.Projects.WUnlockNoPersist()
	}
	return found
}

func (m *Manager) GetProjectsPtfXLS(w io.Writer) {
	m.Projects.RLock()
	defer m.Projects.RUnlock()
	WritePortfolioToXLS(m.Projects.GetProjectsPtf(), w)
}

func (m *Manager) UpdateStat(r io.Reader) {
	m.Stats.UpdateFrom(r)
}

func getProjectKey(p *Model.Project) (string, string) {
	return p.Client, p.Name
}
