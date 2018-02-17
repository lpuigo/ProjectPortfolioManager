package manager

import (
	"fmt"
	"github.com/lpuig/novagile/manager/datamanager"
	ris "github.com/lpuig/novagile/manager/recordindexedset"
	"github.com/lpuig/novagile/model"
	"github.com/xrash/smetrics"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type StatManager struct {
	*datamanager.DataManager
	stat *ris.RecordLinkedIndexedSet
}

func createRISIndexDescs() []ris.IndexDesc {
	res := []ris.IndexDesc{}
	res = append(res, ris.NewIndexDesc("PrjKey", "CLIENT!PROJECT"))
	res = append(res, ris.NewIndexDesc("Issue", "ISSUE"))
	res = append(res, ris.NewIndexDesc("Summary", "SUMMARY"))
	return res
}

func createRISLinkDescs() []ris.LinkDesc {
	res := []ris.LinkDesc{}
	res = append(res, ris.NewLinkDesc("issue-prj", "Issue", "PrjKey"))
	res = append(res, ris.NewLinkDesc("issue-summary", "Issue", "Summary"))
	return res
}

func newStatSetFrom(r io.Reader) (*ris.RecordLinkedIndexedSet, error) {
	cs := ris.NewRecordLinkedIndexedSet(createRISIndexDescs()...)
	for _, ld := range createRISLinkDescs() {
		err := cs.AddLink(ld)
		if err != nil {
			return nil, fmt.Errorf("newStatSetFrom: %s", err.Error())
		}
	}
	err := cs.AddCSVDataFrom(r)
	if err != nil {
		return nil, fmt.Errorf("newStatSetFrom: %s", err.Error())
	}
	return cs, nil
}

func InitStatManagerFile(file string) error {
	header := "EXTRACT_DATE;PRODUCT;CLIENT!PROJECT;ACTIVITY;ISSUE;INIT_ESTIMATE;TIME_SPENT;REMAIN_TIME;SUMMARY"
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
	sm.DataManager = datamanager.NewDataManager(func() error {
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
	sm.stat, _ = sm.stat.CreateSubSet(createRISIndexDescs(), createRISLinkDescs())
}

func (sm *StatManager) GetStats() *ris.RecordLinkedIndexedSet {
	return sm.stat
}

// GetProjectStatList returns Stats Project List which are not not found in given prjlist
func (sm *StatManager) GetProjectStatList(prjlist map[string]bool) []string {
	res := sm.stat.GetLink("issue-prj").Values()
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
	return sm.stat.GetLink("issue-prj").HasValue(pk)
}

// UpdateFrom updates Stats data with new Stats (CSV Formated) found in r
func (sm *StatManager) UpdateFrom(r io.Reader) (int, error) {
	newStats, err := newStatSetFrom(r)
	if err != nil {
		return 0, fmt.Errorf("UpdateFrom: %s", err.Error())
	}

	SREDesc := ris.NewIndexDesc("SRE", "CLIENT!PROJECT", "TIME_SPENT", "REMAIN_TIME", "INIT_ESTIMATE")
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
		//TODO Update project associated to issue
		oldStats.AddRecord(record)
		added++
	}
	return added, nil
}

// GetProjectSpentWL returns Spent WorkLoad for given project client/name, or error if project, client stat is found
func (sm *StatManager) GetProjectSpentWL(client, name string) (spent float64, err error) {
	_, nums, err := sm.issuesInfosFromProject(client, name)
	if err != nil {
		return 0, err
	}
	colpos, _ := sm.stat.GetRecordColNumByName("EXTRACT_DATE", "TIME_SPENT")
	datePos, spentPos := colpos[0], colpos[1]
	issueDate := map[int]string{}
	issueSpent := map[int]float64{}
	for curIssue, inums := range nums {
		for _, rn := range inums {
			record := sm.stat.GetRecordByNum(rn)
			curDate := record[datePos]
			curWL, _ := strconv.ParseFloat(record[spentPos], 64)
			previousDate, issueFound := issueDate[curIssue]
			if issueFound && curDate < previousDate {
				continue
			}
			issueDate[curIssue] = curDate
			issueSpent[curIssue] = curWL
		}
	}
	spent = 0
	for _, wl := range issueSpent {
		spent += wl
	}
	spent /= 8.0
	return
}

func (sm *StatManager) issuesInfosFromProject(client, name string) (issuesKeys []string, nums [][]int, err error) {
	pk := pjrKey(client, name)
	if !sm.hasStatsForProject(pk) {
		err = fmt.Errorf("No Project Stats for %s/%s", client, name)
		return
	}
	//retrieve all issues associated to prjKey pk
	issuesKeys = sm.stat.GetLink("issue-prj").KeysMatching(pk)
	nums = make([][]int, len(issuesKeys))
	for in, issue := range issuesKeys {
		nums[in] = sm.stat.GetRecordNumsByIndexKey("Issue", issue)
	}
	return
}

// GetProjectStatInfoOnPeriod returns list of issues, dates slices within Given Period, and timeSpent, timeRemaining, timeEstimated doubleslices ([#issue][#date]) for given project client/name
func (sm *StatManager) GetProjectStatInfoOnPeriod(client, name, startDate, endDate string) (issues, summaries []string, sDate string, spent, remaining, estimated [][]float64, err error) {
	issuesKeys, nums, err := sm.issuesInfosFromProject(client, name)
	issues = make([]string, len(issuesKeys))
	summaries = make([]string, len(issuesKeys))
	for i, k := range issuesKeys {
		issues[i] = strings.TrimLeft(k, "!")
		summaries[i] = strings.TrimLeft(sm.stat.GetLink("issue-summary").Get(k, "no summary"), "!")
	}
	// On the result RecordSet
	// Get the available update dates from the project stats
	colpos, _ := sm.stat.GetRecordColNumByName("EXTRACT_DATE", "TIME_SPENT", "REMAIN_TIME", "INIT_ESTIMATE")
	datePos, spentPos, remainingPos, estimatedPos := colpos[0], colpos[1], colpos[2], colpos[3]
	minDate, maxDate := "9999999", "00000000"
	for _, num := range nums {
		for _, n := range num {
			date := sm.stat.GetRecordByNum(n)[datePos]
			if date > maxDate {
				maxDate = date
			}
			if date < minDate {
				minDate = date
			}
		}
	}
	if startDate == "" || startDate > minDate {
		startDate = minDate
	}
	if endDate == "" || endDate < maxDate {
		endDate = maxDate
	}
	sDate = startDate
	// Create Date list (chronologically sorted from start-end dates) => result dates slice
	var dates []string
	dates, err = dateSlice(startDate, endDate)
	if err != nil {
		return
	}
	// create result S, R, E slice with Date List length
	initDS(&spent, len(issues), len(dates))
	initDS(&remaining, len(issues), len(dates))
	initDS(&estimated, len(issues), len(dates))

	for ii, ik := range issuesKeys {
		irs, errirs := sm.stat.CreateSubSet(
			[]ris.IndexDesc{
				ris.NewIndexDesc("Dates", "EXTRACT_DATE"),
				ris.NewIndexDesc("IssueDate", "ISSUE", "EXTRACT_DATE"),
			},
			nil,
		)
		if errirs != nil {
			err = errirs
			return
		}
		irs.AddRecords(sm.stat.GetRecordsByNums(nums[ii]))
		dateKeys := irs.GetIndexKeys("Dates")
		sort.Strings(dateKeys)
		for di, dk := range dates {
			dateKey := "!" + dk
			r := irs.GetRecordsByIndexKey("IssueDate", ik+dateKey)
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
						r := irs.GetRecordsByIndexKey("IssueDate", ik+dateKeys[i-1])
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
	d1, err := model.DateFromJSString(startDate)
	if err != nil {
		return nil, fmt.Errorf("misformated startDate '%s'", startDate)
	}
	d2, err := model.DateFromJSString(endDate)
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
