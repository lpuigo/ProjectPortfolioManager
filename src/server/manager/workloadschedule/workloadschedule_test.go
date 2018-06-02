package workloadschedule

import (
	"encoding/json"
	"github.com/lpuig/prjptf/src/server/model"
	"strings"
	"testing"
)

func TestCalcWeeks(t *testing.T) {
	b, e, w := calcWeeks(1, 2)

	curw := model.Today().GetMonday().String()

	if len(w) != 4 {
		t.Error("calcWeeks returns unexpected number of value (4 expected), instead got", w)
	}
	if w[1] != curw {
		t.Errorf("calcWeeks returns unexpected value in week list %v (second should be %s)", w, curw)
	}

	bD, _ := model.DateFromJSString(b)
	eD, _ := model.DateFromJSString(e)
	eD = eD.AddDays(-7)
	bW := bD.String()
	eW := eD.String()

	if w[0] != bW || w[len(w)-1] != eW {
		t.Errorf("calcWeeks returns inconsistent values: begWeek(=%s) <> w[0](=%s) or endWeek(=%s) <> w[3](=%s)", bW, w[0], eW, w[3])
	}
}

func TestWeekCoverage(t *testing.T) {
	monday := "2018-05-07"
	friday := "2018-05-11"

	for _, s := range []struct {
		pBeg    string
		pEnd    string
		expect  float64
		comment string
	}{
		{"2018-01-01", "2018-05-06", 0, "full before"},
		{"2018-05-14", "2018-05-30", 0, "full after"},
		{"2018-01-01", "2018-05-30", 5, "much longer"},
		{monday, friday, 5, "exactly the week"},
		{"2018-01-01", friday, 5, "ends at end of week"},
		{monday, "2018-05-30", 5, "starts at begin of week"},
		{"2018-01-01", "2018-05-08", 2, "ends before week end"},
		{"2018-05-09", "2018-05-30", 3, "starts before week start"},
		{"2018-05-09", "2018-05-10", 2, "all within the week"},
		{"2018-05-08", "2018-05-08", 1, "one day project"},
		{monday, monday, 1, "one day Monday project"},
		{friday, friday, 1, "one day Friday project"},
	} {
		res := calcWeekCoverage(monday, friday, s.pBeg, s.pEnd)
		if res != s.expect {
			t.Errorf("case '%s': calcWeekCoverage(%s, %s, %s, %s) returns %f instead of %f", s.comment, monday, friday, s.pBeg, s.pEnd, res, s.expect)
		}
	}
}

func createPrjs(t *testing.T) []*model.Project {
	const testData = `
[
{	"id": 0,
	"client": "Way Before",	
	"name": "KO", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
				"Pilot End": "2017-09-01",
				"RollOut": "2017-09-05"
		}}]}
},
{	"id": 1,
	"client": "Way After",	
	"name": "KO", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
				"Pilot End": "2019-09-01",
				"RollOut": "2019-09-05"
		}}]}
},
{	"id": 2,
	"client": "All Around",	
	"name": "OK", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
				"KickOff": "2018-01-01",
				"RollOut": "2018-12-31"
		}}]}
},
{	"id": 3,
	"client": "Ends before time span",	
	"name": "OK", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
				"KickOff": "2018-01-01",
				"RollOut": "2018-05-06"
		}}]}
},
{	"id": 4,
	"client": "Begin after time span",	
	"name": "OK", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
				"KickOff": "2018-05-07",
				"RollOut": "2018-12-31"
		}}]}
},
{	"id": 5,
	"client": "No Date",	
	"name": "KO", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
		}}]}
},
{	"id": 6,
	"client": "less than one week",	
	"name": "OK", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
				"KickOff": "2018-05-09",
				"RollOut": "2018-05-10"
		}}]}
},
{	"id": 7,
	"client": "Beg/Ends at mid week",	
	"name": "OK", "lead_dev": "Laurent", "status": "6 - Done", "type": "Acti", 
	"forecast_wl": 1, "current_wl": 0,
	"situation": {
		"stds": [{
			"update": "2017-11-26",
			"milestones": {
				"KickOff": "2018-05-03",
				"RollOut": "2018-05-15"
		}}]}
}
]
`

	res := []*model.Project{}
	err := json.NewDecoder(strings.NewReader(testData)).Decode(&res)
	if err != nil {
		t.Fatal("could not unmarshal testData:", err.Error())
	}
	return res
}

func TestCalcWorkloadSchedule(t *testing.T) {
	prjs := createPrjs(t)
	res := Calc(prjs)

	if len(res.Weeks) != 13 {
		t.Error("Weeks is has not 13 values:%v", res.Weeks)
	}

	for _, ws := range res.Records {
		if prjs[ws.Id].Name == "KO" {
			t.Error("unexpected WS records found for: %s", prjs[ws.Id].Client)
		}
		t.Logf("Worload Schedule for '%s' (id %d); %v", prjs[ws.Id].Client, ws.Id, ws.WorkLoads)
	}
}
