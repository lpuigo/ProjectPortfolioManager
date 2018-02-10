package model

import (
	"encoding/json"
	"testing"
)

const ps1JSONSting = `{"id":1,"dates":"2017-11-27","values":{"Estimated":{"0":5,"1":5},"Remaining":{"0":5,"1":3},"Spent":{"0":0,"1":2}}}`

func TestNewProjectStat(t *testing.T) {
	ps := NewProjectStat()
	ps.Id = 1

	d, _ := DateFromJSString("2017-11-27")
	ps.StartDate = d
	ps.AddValues(d, 0, 5, 5)
	ps.AddValues(d.AddDays(1), 2, 3, 5)

	psjson, err := json.Marshal(ps)
	if err != nil {
		t.Errorf("ProjectStat Marshal returns: %s", err.Error())
	}
	if string(psjson) != ps1JSONSting {
		t.Errorf("ProjectStat Marshal returns wrong result : \n'%s'\ninstead of\n'%s'", string(psjson), ps1JSONSting)
	}

	println(ps.String())
}
