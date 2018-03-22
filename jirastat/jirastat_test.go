package jirastat

import (
	"fmt"
	"testing"
)

const (
	testFile = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Stat_Tempo\extract Worklog 2018-03-13.csv`
)

func displayJs(js *JiraStat) {
	for _, r := range js.Stats.GetRecords() {
		fmt.Println(r)
	}
}

func TestJiraStat_SpentHourBy(t *testing.T) {
	js := NewJiraStat()

	if err := js.LoadFromFile(testFile); err != nil {
		t.Fatal("js.LoadFromFile returns", err.Error())
	}

	keys, values, err := js.SpentHourBy("LotClient")
	if err != nil {
		t.Fatal("JiraStat.SpentHourBy returns", err.Error())
	}

	for i, k := range keys {
		fmt.Printf("%s : %0.3f\n", k, values[i]/8)
	}

}
