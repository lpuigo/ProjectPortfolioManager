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

	sm, err := NewStatManagerFromPersistFile(statfile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve Stat Portfolio Data: %s", err.Error()))
	}
	m.Stats = sm

	return m, nil
}

func (m *Manager) GetPrjPtf(w io.Writer) {
	prjs := []*fm.Project{}
	// TODO figure out how to avoid allocating full ptf clone before marshaling
	m.Projects.RLock() // Ensure PTf is not modified while being cloned
	for _, p := range m.Projects.GetProjectsPtf().Projects {
		prjs = append(prjs, fm.CloneBEProject(p))
	}
	m.Projects.RUnlock()
	json.NewEncoder(w).Encode(prjs)
}

func (m *Manager) GetPrjById(id int) *Model.Project {
	return m.Projects.GetProjectsPtf().GetPrjById(id)
}

func (m *Manager) UpdateProject(op, np *Model.Project) {
	m.Projects.WLock()
	defer m.Projects.WUnlock()
	op.Update(np)
}

func (m *Manager) CreateProject(p *Model.Project) *Model.Project {
	m.Projects.WLock()
	defer m.Projects.WUnlock()
	m.Projects.GetProjectsPtf().AddPrj(p)
	return p
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

func (m *Manager) UpdateStatFromCSVFile(csvfile string) error {
	//m.Projects.RLock()
	//m.Stats.WLock()
	////TODO create new StatsPtf and then, if no error nor warning, append it to actual Manager.Stats
	//num, err, warns := UpdateStatPortfolioFromCSVFile(csvfile, m.Projects.GetProjectsPtf(), m.Stats.GetStatsPtf())
	//m.Projects.RUnlock()
	//if err != nil {
	//	m.Stats.WUnlockNoPersist()
	//	return err
	//}
	//log.Printf("Statfile %s processed : %d stats added\n", csvfile, num)
	//if warns.HasWarnings() {
	//	log.Printf("Warnings :\n%s", warns.Warning(""))
	//}
	//if num == 0 {
	//	m.Stats.WUnlockNoPersist()
	//} else {
	//	m.Stats.WUnlock()
	//}
	return nil
}
