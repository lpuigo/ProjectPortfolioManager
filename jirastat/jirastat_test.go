package jirastat

import (
	"fmt"
	"testing"
)

const (
	testFile = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Stat_Tempo\extract Worklog 2018-03-13.csv`
)

func displayJs(js *JiraStat) {
	for _, r := range js.GetRecords() {
		fmt.Println(r)
	}
}

func TestJiraStat_LoadFromFile(t *testing.T) {
	js := NewJiraStat()

	if err := js.LoadFromFile(testFile); err != nil {
		t.Fatal("js.LoadFromFile returns", err.Error())
	}

	displayJs(js)

}
