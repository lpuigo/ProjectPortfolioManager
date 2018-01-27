package Manager

import (
	"fmt"
	"github.com/lpuig/Novagile/Manager/DataManager"
	ris "github.com/lpuig/Novagile/Manager/RecordIndexedSet"
	"github.com/lpuig/Novagile/Model"
	"github.com/xrash/smetrics"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type StatManager struct {
	*DataManager.DataManager
	stat *ris.RecordIndexedSet
}

func createRISIndexDescs() []ris.IndexDesc {
	res := []ris.IndexDesc{}
	res = append(res, ris.NewIndexDesc("PrjKey", "CLIENT!PROJECT"))
	res = append(res, ris.NewIndexDesc("Issue", "ISSUE"))
	return res
}

func newStatSetFrom(r io.Reader) (*ris.RecordIndexedSet, error) {
	cs := ris.NewRecordIndexedSet(createRISIndexDescs()...)
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

func (sm *StatManager) ClearStats() {
	sm.stat, _ = sm.stat.CreateSubRecordIndexedSet(createRISIndexDescs()...)
}

func (sm *StatManager) GetStats() *ris.RecordIndexedSet {
	return sm.stat
}

func (sm *StatManager) GetProjectStatList(prjlist map[string]bool) []string {
	res := sm.stat.GetIndexKeys("PrjKey")
	for i, s := range res {
		if _, exist := prjlist[s]; exist {
			res[i] = ""
			continue
		}
		res[i] = strings.TrimLeft(s, "!")
	}
	sort.Strings(res)
	// remove leading "" project signature
	found := 0
	for _, s := range res {
		if s == "" {
			found++
		} else {
			break
		}
	}
	return res[found:]
}

type slicePair struct {
	mainSlice []string
	distSlice []float64
}

func (sbd slicePair) Len() int {
	return len(sbd.mainSlice)
}

func (sbd slicePair) Swap(i, j int) {
	sbd.mainSlice[i], sbd.mainSlice[j] = sbd.mainSlice[j], sbd.mainSlice[i]
	sbd.distSlice[i], sbd.distSlice[j] = sbd.distSlice[j], sbd.distSlice[i]
}

func (sbd slicePair) Less(i, j int) bool {
	return sbd.distSlice[j] < sbd.distSlice[i]
}

func SortByDist(list []string, dist []float64) {
	sbd := slicePair{mainSlice: list, distSlice: dist}
	sort.Sort(sbd)
}

func (sm *StatManager) GetProjectStatListSortedBySimilarity(signature string, prjlist map[string]bool) []string {
	list := sm.GetProjectStatList(prjlist)
	dist := make([]float64, len(list))
	compareString := strings.ToUpper(signature)
	for i, s := range list {
		dist[i] = smetrics.JaroWinkler(compareString, strings.ToUpper(s), 0.7, 1)
	}
	SortByDist(list, dist)
	return list
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

// UpdateFrom updates Stats data with new Stats (CSV Formated) found in r
func (sm *StatManager) UpdateFrom(r io.Reader) (int, error) {
	newStats, err := newStatSetFrom(r)
	if err != nil {
		return 0, fmt.Errorf("UpdateFrom: %s", err.Error())
	}

	SREDesc := ris.NewIndexDesc("SRE", "TIME_SPENT", "REMAIN_TIME", "INIT_ESTIMATE")
	dateDesc := ris.NewIndexDesc("Date", "EXTRACT_DATE")
	oldStats := sm.GetStats()

	oldStatSREKey, err := oldStats.GetKeyGeneratorByIndexDesc(SREDesc)
	if err != nil {
		return 0, err
	}
	newStatSREKey, err := newStats.GetKeyGeneratorByIndexDesc(SREDesc)
	if err != nil {
		return 0, err
	}

	added := 0
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
	return added, nil
}

// GetProjectSpentWL returns Spent WorkLoad for given project client/name, or error if project, client stat is found
func (sm *StatManager) GetProjectSpentWL(client, name string) (spent float64, err error) {
	pk := pjrKey(client, name)
	if !sm.hasStatsForProject(pk) {
		err = fmt.Errorf("No Project Stats for %s/%s", client, name)
		return
	}
	//retrieve all issues associated to prjKey pk
	ps, errss := sm.stat.CreateSubRecordIndexedSet(
	//ris.NewIndexDesc("IssueDate", "ISSUE", "EXTRACT_DATE"),
	)
	if errss != nil {
		err = fmt.Errorf("PrjSubSet: %s", errss.Error())
	}
	ps.AddRecords(sm.stat.GetRecordsByIndexKey("PrjKey", pk))

	colpos, _ := ps.GetRecordColNumByName("ISSUE", "EXTRACT_DATE", "TIME_SPENT")
	issuePos, datePos, spentPos := colpos[0], colpos[1], colpos[2]
	issueDate := map[string]string{}
	issueSpent := map[string]float64{}
	for _, record := range ps.GetRecords() {
		curIssue := record[issuePos]
		curDate := record[datePos]
		curWL, _ := strconv.ParseFloat(record[spentPos], 64)
		previousDate, issueFound := issueDate[curIssue]
		if issueFound && curDate < previousDate {
			continue
		}
		issueDate[curIssue] = curDate
		issueSpent[curIssue] = curWL
	}
	spent = 0
	for _, wl := range issueSpent {
		spent += wl
	}
	spent /= 8.0
	return
}

func (sm *StatManager) initStatInfoRecordSet(client, name string) (issuesKeys []string, is *ris.RecordIndexedSet, err error) {
	pk := pjrKey(client, name)
	if !sm.hasStatsForProject(pk) {
		err = fmt.Errorf("No Project Stats for %s/%s", client, name)
		return
	}
	//retrieve all issues associated to prjKey pk
	ps, errss := sm.stat.CreateSubRecordIndexedSet(
		ris.NewIndexDesc("Issue", "ISSUE"),
	)
	if errss != nil {
		err = fmt.Errorf("PrjSubSet: %s", errss.Error())
	}
	ps.AddRecords(sm.stat.GetRecordsByIndexKey("PrjKey", pk))
	// Keep Issue list => result issues slice
	issuesKeys = ps.GetIndexKeys("Issue")
	sort.Strings(issuesKeys)

	//retrieve all issues found in ps
	is, errss = sm.stat.CreateSubRecordIndexedSet(
		ris.NewIndexDesc("Issue", "ISSUE"),
		ris.NewIndexDesc("Dates", "EXTRACT_DATE"),
		ris.NewIndexDesc("IssueDate", "ISSUE", "EXTRACT_DATE"),
	)
	if errss != nil {
		err = fmt.Errorf("IssueSubSet: %s", errss.Error())
	}
	// For each identified issues,
	for _, ik := range issuesKeys {
		// retrieve all record related to this issue in a new SubRecordSet (with Date Index)
		is.AddRecords(sm.stat.GetRecordsByIndexKey("Issue", ik))
	}
	return
}

// GetProjectStatInfo returns list of issues, dates slices, and timeSpent, timeRemaining, timeEstimated doubleslices ([#issue][#date]) for given project client/name
func (sm *StatManager) GetProjectStatInfo(client, name string) (issues, dates []string, spent, remaining, estimated [][]float64, err error) {
	issuesKeys, is, err := sm.initStatInfoRecordSet(client, name)
	issues = make([]string, len(issuesKeys))
	for i, k := range issuesKeys {
		issues[i] = strings.TrimLeft(k, "!")
	}
	// On the result RecordSet
	// Keep Date list (chronologically sorted) => result dates slice
	dateKeys := is.GetIndexKeys("Dates")
	sort.Strings(dateKeys)
	dates = make([]string, len(dateKeys))
	for i, k := range dateKeys {
		dates[i] = strings.TrimLeft(k, "!")
	}
	// create result S, R, E slice with Date List length
	initDS(&spent, len(issues), len(dates))
	initDS(&remaining, len(issues), len(dates))
	initDS(&estimated, len(issues), len(dates))
	colpos, _ := is.GetRecordColNumByName("TIME_SPENT", "REMAIN_TIME", "INIT_ESTIMATE")
	spentPos, remainingPos, estimatedPos := colpos[0], colpos[1], colpos[2]
	// For each Date,
	for ii, ik := range issuesKeys {
		for di, dk := range dateKeys {
			r := is.GetRecordsByIndexKey("IssueDate", ik+dk)
			if r == nil && di > 0 {
				spent[ii][di] = spent[ii][di-1]
				remaining[ii][di] = remaining[ii][di-1]
				estimated[ii][di] = estimated[ii][di-1]
			}
			if r == nil {
				continue
			}
			spent[ii][di], err = stringToWL(r[0][spentPos])
			remaining[ii][di], err = stringToWL(r[0][remainingPos])
			estimated[ii][di], err = stringToWL(r[0][estimatedPos])
		}
	}
	return
}

// GetProjectStatInfoOnPeriod returns list of issues, dates slices within Given Period, and timeSpent, timeRemaining, timeEstimated doubleslices ([#issue][#date]) for given project client/name
func (sm *StatManager) GetProjectStatInfoOnPeriod(client, name, startDate, endDate string) (issues, dates []string, spent, remaining, estimated [][]float64, err error) {
	issuesKeys, is, err := sm.initStatInfoRecordSet(client, name)
	issues = make([]string, len(issuesKeys))
	for i, k := range issuesKeys {
		issues[i] = strings.TrimLeft(k, "!")
	}
	// On the result RecordSet
	// Create Date list (chronologically sorted from start-end dates) => result dates slice
	dates, err = dateSlice(startDate, endDate)
	if err != nil {
		return
	}
	// create result S, R, E slice with Date List length
	initDS(&spent, len(issues), len(dates))
	initDS(&remaining, len(issues), len(dates))
	initDS(&estimated, len(issues), len(dates))
	colpos, _ := is.GetRecordColNumByName("TIME_SPENT", "REMAIN_TIME", "INIT_ESTIMATE")
	spentPos, remainingPos, estimatedPos := colpos[0], colpos[1], colpos[2]

	for ii, ik := range issuesKeys {
		irs, errirs := is.CreateSubRecordIndexedSet(ris.NewIndexDesc("Dates", "EXTRACT_DATE"))
		if errirs != nil {
			err = errirs
			return
		}
		irs.AddRecords(is.GetRecordsByIndexKey("Issue", ik))
		dateKeys := irs.GetIndexKeys("Dates")
		sort.Strings(dateKeys)
		for di, dk := range dates {
			dateKey := "!" + dk
			r := is.GetRecordsByIndexKey("IssueDate", ik+dateKey)
			if r == nil {
				if di > 0 {
					spent[ii][di] = spent[ii][di-1]
					remaining[ii][di] = remaining[ii][di-1]
					estimated[ii][di] = estimated[ii][di-1]
				} else { // first date : init values
					if dateKey < dateKeys[0] {
						//spent[ii][di] = 0.0
						//remaining[ii][di] = 0.0
						//estimated[ii][di] = 0.0
					} else {
						var i int
						for i = 1; i < len(dateKeys); i++ {
							if dateKey < dateKeys[i] {
								break
							}
						}
						r := is.GetRecordsByIndexKey("IssueDate", ik+dateKeys[i-1])
						spent[ii][di], err = stringToWL(r[0][spentPos])
						remaining[ii][di], err = stringToWL(r[0][remainingPos])
						estimated[ii][di], err = stringToWL(r[0][estimatedPos])
					}
				}
				continue
			}
			spent[ii][di], err = stringToWL(r[0][spentPos])
			remaining[ii][di], err = stringToWL(r[0][remainingPos])
			estimated[ii][di], err = stringToWL(r[0][estimatedPos])
		}
	}
	return
}

func initDS(ds *[][]float64, len1, len2 int) {
	*ds = make([][]float64, len1)
	for i, _ := range *ds {
		(*ds)[i] = make([]float64, len2)
	}
}

func stringToWL(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return f / 8.0, nil
}

func dateSlice(startDate, endDate string) ([]string, error) {
	d1, err := Model.DateFromJSString(startDate)
	if err != nil {
		return nil, fmt.Errorf("misformated startDate '%s'", startDate)
	}
	d2, err := Model.DateFromJSString(endDate)
	if err != nil {
		return nil, fmt.Errorf("misformated endDate '%s'", endDate)
	}
	nbdays := d2.DaysSince(d1)
	res := make([]string, nbdays+1)
	res[0] = startDate
	for i := 1; i <= nbdays; i++ {
		res[i] = d1.AddDays(i).StringJS()
	}
	return res, nil
}
