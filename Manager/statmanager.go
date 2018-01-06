package Manager

import (
	"errors"
	"fmt"
	"github.com/lpuig/Novagile/Manager/CSVStats"
	"github.com/lpuig/Novagile/Manager/DataManager"
	"log"
)

type StatManager struct {
	*DataManager.DataManager
	stat *CSVStats.CSVStats
}

func createCSVStatsIndexDescs() []CSVStats.IndexDesc {
	res := []CSVStats.IndexDesc{}
	res = append(res, CSVStats.NewIndexDesc("PrjKey", "CLIENT!PROJECT"))
	res = append(res, CSVStats.NewIndexDesc("Issue", "ISSUE"))
	return res
}

func NewStatManagerFromPersistFile(file string) (*StatManager, error) {
	sm := &StatManager{}
	sm.DataManager = DataManager.NewDataManager(func() error {
		return sm.stat.WriteCSVToFile(file)
	})
	sm.stat = CSVStats.NewCSVStats(createCSVStatsIndexDescs()...)
	err := sm.stat.AddCSVDataFromFile(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("NewStatManagerFromPersistFile '%s' : %s", file, err.Error()))
	}
	log.Printf("Stats loaded (%d projects, %d records(s))\n", len(sm.stat.GetKeys("PrjKey")), sm.stat.Len())
	return sm, nil
}

func (sm *StatManager) GetStats() *CSVStats.CSVStats {
	return sm.stat
}
