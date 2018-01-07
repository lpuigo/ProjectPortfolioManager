package Manager

import (
	"fmt"
	"os"
	"sort"
	"testing"
	"time"
)

const StatFile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-03.csv`
const StatFile2 = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-04.csv`
const StatFile0 = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\test_extract_init.csv`

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
	sm, err := NewStatManagerFromFile(StatFile)
	if err != nil {
		t.Fatalf("NewStatManagerFromFile: %s", err.Error())
	}

	f, err := os.Open(StatFile2)
	if err != nil {
		t.Fatalf("StatManager_UpdateFrom: %s", err.Error())
	}
	defer f.Close()
	err = sm.UpdateFrom(f)
	if err != nil {
		t.Fatalf("UpdateFrom returns %s", err.Error())
	}
	time.Sleep(4 * time.Second)
}

func TestInitActualDataWithProd(t *testing.T) {
	UpdateStatPortfolioFromCSVFile()
}
