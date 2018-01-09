package Manager

import (
	"fmt"
	"github.com/lpuig/Novagile/Manager/DataManager"
	ris "github.com/lpuig/Novagile/Manager/RecordIndexedSet"
	"io"
	"log"
	"os"
	"sort"
)

type StatManager struct {
	*DataManager.DataManager
	stat *ris.RecordIndexedSet
}

func createCSVStatsIndexDescs() []ris.IndexDesc {
	res := []ris.IndexDesc{}
	res = append(res, ris.NewIndexDesc("PrjKey", "CLIENT!PROJECT"))
	res = append(res, ris.NewIndexDesc("Issue", "ISSUE"))
	return res
}

func newStatSetFrom(r io.Reader) (*ris.RecordIndexedSet, error) {
	cs := ris.NewRecordIndexedSet(createCSVStatsIndexDescs()...)
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

func (sm *StatManager) GetStats() *ris.RecordIndexedSet {
	return sm.stat
}

func (sm *StatManager) GetProjectStatList() []string {
	res := sm.stat.GetIndexKeys("PrjKey")
	sort.Strings(res)
	return res
}

func (sm *StatManager) HasStatsForProject(client, name string) bool {
	return sm.hasStatsForProject(pjrKey(client, name))
}

func pjrKey(client, name string) string {
	return "!" + client + "!" + name
}

func (sm *StatManager) hasStatsForProject(pk string) bool {
	return sm.stat.HasIndexKey("PrjKey", pk)
}

// UpdateFrom updates Stats data with new Stats (CSV Formated) found in r (StatManager is WriteLocked during process)
func (sm *StatManager) UpdateFrom(r io.Reader) error {
	newStats, err := newStatSetFrom(r)
	if err != nil {
		return fmt.Errorf("UpdateFrom: %s", err.Error())
	}

	SREDesc := ris.NewIndexDesc("SRE", "TIME_SPENT", "REMAIN_TIME", "INIT_ESTIMATE")
	dateDesc := ris.NewIndexDesc("Date", "EXTRACT_DATE")
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

// GetProjectStatInfo returns list of issues, dates, and timeSpent, timeRemaining, timeEstimated slices for given project client/name
func (sm *StatManager) GetProjectStatInfo(client, name string) (issues, dates []string, spent, remaining, estimated [][]float64, err error) {
	pk := pjrKey(client, name)
	if !sm.hasStatsForProject(pk) {
		err = fmt.Errorf("No Project Stats for %s/%s", client, name)
		return
	}
	//retrieve all issues associated to prjKey pk
	ps, err := sm.stat.CreateSubRecordIndexedSet(
		ris.NewIndexDesc("Issue", "ISSUE"),
	)
	for _, r := range sm.stat.GetRecordsByIndexKey("PrjKey", pk) {
		ps.AddRecord(r)
	}
	// For each identified issues,
	// retrieve all record related to this issue in a new SubRecordSet (with Date Index)

	// On the result RecordSet
	// Keep Issue list => result issues slice
	// Keep Date list (chronologically sorted) => result dates slice
	// create result S, R, E slice with Date List length
	// For each Date,
	// store Date, Sum of S, R, E
	//TODO Implement!!
	return
}
