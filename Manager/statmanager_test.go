package Manager

import (
	"fmt"
	"sort"
	"testing"
)

const StatFile = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\export Jira\extract 2018-01-03.csv`

func TestNewStatManagerFromPersistFile(t *testing.T) {
	sm, err := NewStatManagerFromPersistFile(StatFile)
	if err != nil {
		t.Errorf("NewStatManagerFromPersistFile: %s", err.Error())
	}

	prjs := sm.stat.GetKeys("PrjKey")
	sort.Strings(prjs)
	fmt.Printf("Known project: %s\n", prjs)

	issues := sm.stat.GetKeys("Issue")
	sort.Strings(issues)
	fmt.Printf("Known Jira: %s\n", issues)
}
