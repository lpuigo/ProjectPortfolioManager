package Manager

import (
	"github.com/lpuig/Novagile/Manager/DataManager"
	"log"
)

type PrjManager struct {
	*DataManager.DataManager
	ptf *PrjPortfolio
}

func NewPrjManagerFromPersistFile(file string) (*PrjManager, error) {
	pm := &PrjManager{}
	pm.DataManager = DataManager.NewDataManager(func() error {
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
