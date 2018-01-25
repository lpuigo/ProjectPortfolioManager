package Manager

import (
	"fmt"
	"github.com/lpuig/Novagile/Manager/DataManager"
	"github.com/lpuig/Novagile/Manager/RecordSet"
	"os"
	"sort"
	"strings"
	"testing"
)

const (
	StatFile  = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-03.csv`
	StatFile2 = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-04.csv`
	StatFile0 = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\test_extract_init.csv`
)

func TestInitStatManagerFile(t *testing.T) {
	sm, err := NewStatManagerFromFile(StatFile0)
	if err != nil {
		t.Fatalf("NewStatManagerFromFile on new file returns '%s'", err.Error())
	}
	defer os.Remove(StatFile0)
	if sm.GetStats().Len() != 0 {
		t.Errorf("NewStatManagerFromFile on new file is not empty")
	}
}

func TestNewStatManagerFromFile(t *testing.T) {
	sm, err := NewStatManagerFromFile(StatFile)
	if err != nil {
		t.Fatalf("NewStatManagerFromFile: %s", err.Error())
	}

	prjs := sm.stat.GetIndexKeys("PrjKey")
	sort.Strings(prjs)
	fmt.Printf("Known project: %s\n", prjs)

	issues := sm.stat.GetIndexKeys("Issue")
	sort.Strings(issues)
	fmt.Printf("Known Jira: %s\n", issues)
}

func TestStatManager_GetStats(t *testing.T) {
	sm, err := NewStatManagerFromFile(StatFile)
	if err != nil {
		t.Fatalf("NewStatManagerFromFile: %s", err.Error())
	}

	s := sm.GetStats()

	prjs := s.GetIndexKeys("PrjKey")
	sort.Strings(prjs)
	fmt.Printf("Known project: %s\n", prjs)
}

func TestStatManager_UpdateFrom(t *testing.T) {
	sm, err := NewStatManagerFromFile(StatFile0)
	if err != nil {
		t.Fatalf("NewStatManagerFromFile: %s", err.Error())
	}
	sm.ClearStats()

	f, err := os.Open(StatFile)
	if err != nil {
		t.Fatalf("StatManager_UpdateFrom: %s", err.Error())
	}
	defer f.Close()
	added, err := sm.UpdateFrom(f)
	if err != nil {
		t.Fatalf("UpdateFrom returns %s", err.Error())
	}
	if added != 155 {
		t.Fatalf("UpdateFrom returns %d added rows, 155 expected", added)
	}
}

func TestStatManager_GetProjectStatList(t *testing.T) {
	sm, err := NewStatManagerFromFile(PrdStatFile)
	if err != nil {
		t.Fatalf("NewStatManagerFromFile: %s", err.Error())
	}
	knownProjects := map[string]bool{"!Bouygues Telecom!DME TVE": true}
	for i, s := range sm.GetProjectStatList(knownProjects) {
		//ls := strings.Split(s, "!")
		fmt.Printf("%2d : %s\n", i, s)
	}
}

func TestStatManager_GetProjectStatListSortedBySimilarity(t *testing.T) {
	sm, err := NewStatManagerFromFile(PrdStatFile)
	if err != nil {
		t.Fatalf("NewStatManagerFromFile: %s", err.Error())
	}
	//knownProjects := map[string]bool{"!Bouygues Telecom!DME TVE": true}
	knownProjects := map[string]bool{}
	project := "Bouygues Telecom!DME TVE"
	statsProjects := sm.GetProjectStatListSortedBySimilarity(project, knownProjects)
	for i, s := range statsProjects {
		fmt.Printf("%2d : %s\n", i, s)
	}
}

func equals(a, b [][]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, fs := range a {
		if len(fs) != len(b[i]) {
			return false
		}
		for j, f := range fs {
			if f != b[i][j] {
				return false
			}
		}
	}
	return true
}

func createTestSM(t *testing.T) *StatManager {
	smSource := `EXTRACT_DATE;PRODUCT;CLIENT!PROJECT;ACTIVITY;ISSUE;INIT_ESTIMATE;TIME_SPENT;REMAIN_TIME
2017-01-01;TestProduct;SomeClient!OtherProject;;Issue0;8.00;0.00;8.00
2017-01-01;TestProduct;;;Issue1;40.00;0.00;40.00
2017-01-02;TestProduct;SomeClient!TestProject;;Issue1;40.00;40.00;0.00
2017-01-02;TestProduct;SomeClient!TESTProject;;Issue2;16.00;8.00;8.00
2017-01-03;TestProduct;SomeClient!TestProject;;Issue2;16.00;16.00;0.00
`
	sm := &StatManager{}
	sm.DataManager = DataManager.NewDataManager(func() error { return nil })
	cs, err := newStatSetFrom(strings.NewReader(smSource))
	if err != nil {
		t.Fatalf("newStatSetFrom returns %s", err.Error())
	}
	sm.stat = cs
	return sm
}

func TestStatManager_GetProjectStatInfo(t *testing.T) {
	sm := createTestSM(t)

	issues, dates, spent, remaining, estimated, err := sm.GetProjectStatInfo("SomeClient", "TestProject")
	if err != nil {
		t.Fatalf("GetProjectStatInfo returns %s", err.Error())
	}
	if !RecordSet.Record(issues).Equals(RecordSet.Record{"Issue1", "Issue2"}) {
		t.Errorf("issues: %s", issues)
	}
	if !RecordSet.Record(dates).Equals(RecordSet.Record{"2017-01-01", "2017-01-02", "2017-01-03"}) {
		t.Errorf("dates: %s", dates)
	}
	if !equals(spent, [][]float64{[]float64{0.0, 5.0, 5.0}, []float64{0.0, 1.0, 2.0}}) {
		t.Errorf("spent %f", spent)
	}
	if !equals(remaining, [][]float64{[]float64{5.0, 0.0, 0.0}, []float64{0.0, 1.0, 0.0}}) {
		t.Errorf("remaining %f", remaining)
	}
	if !equals(estimated, [][]float64{[]float64{5.0, 5.0, 5.0}, []float64{0.0, 2.0, 2.0}}) {
		t.Errorf("estimated: %f", estimated)
	}
}

func TestStatManager_GetProjectSpentWL(t *testing.T) {
	sm := createTestSM(t)
	wl, err := sm.GetProjectSpentWL("SomeClient", "TestProject")
	if err != nil {
		t.Fatalf("GetProjectSpentWL returns %s", err.Error())
	}
	if wl != 7 {
		t.Errorf("GetProjectSpentWL returns unespected value %f instead of 7", wl)
	}
}

func TestStatManager_dateSlice(t *testing.T) {
	tset := []struct {
		startd string
		endd   string
		expres string
	}{
		{"2016-12-31", "2016-12-31", "2016-12-31"},
		{"2016-12-31", "2017-01-01", "2016-12-31 2017-01-01"},
		{"2016-12-31", "2017-01-02", "2016-12-31 2017-01-01 2017-01-02"},
	}

	for _, e := range tset {
		res := dateSlice(e.startd, e.endd)
		if strings.Join(res, " ") != e.expres {
			t.Errorf("dateSlice('%s', '%s') returns %s instead of [%s]", res, e.expres)
		}
	}
}
