package Manager

import (
	"errors"
	"fmt"
	"github.com/lpuig/Novagile/Manager/DataManager"
	"log"
)

type StatManager struct {
	*DataManager.DataManager
	stat *StatPortfolio
}

func NewStatManagerFromPersistFile(file string) (*StatManager, error) {
	sm := &StatManager{}
	sm.DataManager = DataManager.NewDataManager(func() error {
		return sm.GetStatsPtf().WriteJsonFile(file)
	})
	sp, err := NewStatPortfolioFromJSONFile(file)
	if err != nil {
		log.Println()
		return nil, errors.New(fmt.Sprintf("Unable to load StatPortfolio from file '%s' : %s", file, err.Error()))
	}
	if sp != nil {
		log.Printf("StatPortfolio loaded (%d project(s))\n", len(sp.Stats))
	}
	sm.stat = sp
	return sm, nil
}

func InitStatManagerPersistFile(file string) error {
	sm := &StatManager{}
	sm.DataManager = DataManager.NewDataManager(nil)
	sm.stat = NewStatPortfolio()
	return sm.GetStatsPtf().WriteJsonFile(file)
}

func (sm *StatManager) GetStatsPtf() *StatPortfolio {
	return sm.stat
}
