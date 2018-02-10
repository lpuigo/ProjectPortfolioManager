package manager

import (
	"github.com/lpuig/novagile/manager/datamanager"
	"log"
)

type PrjManager struct {
	*datamanager.DataManager
	ptf *PrjPortfolio
}

func NewPrjManagerFromPersistFile(file string) (*PrjManager, error) {
	pm := &PrjManager{}
	pm.DataManager = datamanager.NewDataManager(func() error {
		return pm.GetProjectsPtf().WriteJsonFile(file)
	})
	p, err := NewPrjPortfolioFromJSONFile(file)
	if err != nil {
		log.Println()
		return nil, err
	}
	if p != nil {
		log.Printf("PrjPortfolio loaded (%d project(s))\n", len(p.Projects))
	}
	pm.ptf = p
	return pm, nil
}

func (pm *PrjManager) GetProjectsPtf() *PrjPortfolio {
	return pm.ptf
}
