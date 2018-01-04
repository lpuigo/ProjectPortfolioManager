package Model

import (
	"encoding/json"
	"testing"
)

const (
	PRJ1Json   = `{"id":1,"client":"a client","name":"prj name","lead_dev":"Laurent","status":"3 - Dev","type":"Novagile","forecast_wl":20,"current_wl":2.2,"comment":"all clear","situation":` + SIT1Json + `}`
	PRJ1String = `Project :
	Id : 1
	Client : a client
	Name : prj name
	Lead Dev : Laurent
	Status : 3 - Dev
	Type : Novagile
	Forecast WorkLoad : 20.0
	Current WorkLoad : 2.2
	Comment : all clear
	Situation : ` + SIT1String + `
.
`
)

func makePRJ1() *Project {
	p := NewProject()
	p.Id = 1
	p.Client = "a client"
	p.Name = "prj name"
	p.Status = StatutDev
	p.Type = TypoNovagile
	p.LeadDev = "Laurent"
	p.Comment = "all clear"
	p.ForecastWL = 20.0
	p.CurrentWL = 2.2
	p.Situation = makeSIT1()
	return p
}

func TestProject_String(t *testing.T) {
	prj := makePRJ1()
	sprj := prj.String()
	if sprj != PRJ1String {
		t.Error("Project.String() returns improper value", sprj)
	}
}

func TestProject_Marshal(t *testing.T) {
	prj := makePRJ1()
	b, err := json.Marshal(prj)
	if err != nil {
		t.Error("Project Marshal returns error", err)
	}
	sb := string(b)
	if sb != PRJ1Json {
		t.Error("Project Marshal returns improper value", sb)
	}
}
