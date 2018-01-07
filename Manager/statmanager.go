package Manager

import (
	"fmt"
	"github.com/lpuig/Novagile/Manager/CSVStats"
	"github.com/lpuig/Novagile/Manager/DataManager"
	"io"
	"log"
	"os"
	"sort"
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

func newStatSetFrom(r io.Reader) (*CSVStats.CSVStats, error) {
	cs := CSVStats.NewCSVStats(createCSVStatsIndexDescs()...)
	err := cs.AddCSVDataFrom(r)
	if err != nil {
		return nil, fmt.Errorf("newStatSetFrom: %s", err.Error())
	}
	return cs, nil
}

func InitStatManagerFile(file string) error {
	header := "EXTRACT_DATE;PRODUCT;CLIENT!PROJECT;ACTIVITY;ISSUE;INIT_ESTIMATE;TIME_SPENT;REMAIN_TIME"
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(([]byte)(header))
	if err != nil {
		return err
	}
	return nil
}

func NewStatManagerFromFile(file string) (*StatManager, error) {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			err := InitStatManagerFile(file)
			if err != nil {
				return nil, err
			}
		}
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("File '%s' : %s", file, err.Error())
	}
	defer f.Close()

	sm := &StatManager{}
	sm.DataManager = DataManager.NewDataManager(func() error {
		return sm.stat.WriteCSVToFile(file)
	})
	cs, err := newStatSetFrom(f)
	if err != nil {
		return nil, fmt.Errorf("File '%s' : %s", file, err.Error())
	}
	sm.stat = cs
	log.Printf("Stats loaded (%d projects, %d records(s))\n", len(sm.stat.GetIndexKeys("PrjKey")), sm.stat.Len())
	return sm, nil
}

func (sm *StatManager) GetStats() *CSVStats.CSVStats {
	return sm.stat
}

func (sm *StatManager) GetProjectStatList() []string {
	res := sm.stat.GetIndexKeys("PrjKey")
	sort.Strings(res)
	return res
}

func (sm *StatManager) HasStatsForProject(client, name string) bool {
	pk := "!" + client + "!" + name
	return sm.stat.HasIndexKey("PrjKey", pk)
}

// UpdateFrom updates Stats data with new Stats (CSV Formated) found in r (StatManager is WriteLocked during process)
func (sm *StatManager) UpdateFrom(r io.Reader) error {
	newStats, err := newStatSetFrom(r)
	if err != nil {
		return fmt.Errorf("UpdateFrom: %s", err.Error())
	}

	SREDesc := CSVStats.NewIndexDesc("SRE", "TIME_SPENT", "REMAIN_TIME", "INIT_ESTIMATE")
	dateDesc := CSVStats.NewIndexDesc("Date", "EXTRACT_DATE")
	oldStats := sm.GetStats()

	oldStatSREKey, err := oldStats.GetKeyGeneratorByIndexDesc(SREDesc)
	if err != nil {
		return err
	}
	newStatSREKey, err := newStats.GetKeyGeneratorByIndexDesc(SREDesc)
	if err != nil {
		return err
	}

	added := 0
	sm.WLock()
	for _, record := range newStats.GetRecords() {
		issueKey := newStats.GetRecordKeyByIndex("Issue", record)
		if oldStats.HasIndexKey("Issue", issueKey) {
			if newStatSREKey(record) == oldStatSREKey(oldStats.Max("Issue", issueKey, dateDesc)) {
				continue
			}
		}
		oldStats.AddRecord(record)
		added++
	}
	if added == 0 {
		sm.WUnlock()
	} else {
		sm.WUnlockWithPersist()
	}
	return nil
}
