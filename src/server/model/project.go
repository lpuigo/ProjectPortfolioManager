package model

import (
	"strconv"
)

type Project struct {
	Id         int        `json:"id"`
	Client     string     `json:"client"`
	Name       string     `json:"name"`
	Risk       int        `json:"risk,omitempty"`
	LeadPS     string     `json:"lead_ps,omitempty"`
	LeadDev    string     `json:"lead_dev,omitempty"`
	Status     string     `json:"status"`
	Type       string     `json:"type"`
	ForecastWL float64    `json:"forecast_wl"`
	CurrentWL  float64    `json:"current_wl"`
	Comment    string     `json:"comment,omitempty"`
	Situation  Situations `json:"situation"`
}

func NewProject() *Project {
	p := &Project{}
	p.Id = 0
	p.Client = ""
	p.Name = ""
	p.Risk = 0
	p.LeadPS = ""
	p.LeadDev = ""
	p.Status = ""
	p.Type = ""
	p.ForecastWL = 0
	p.CurrentWL = 0
	p.Comment = ""
	p.Situation = NewSituations()
	return p
}

func (p Project) String() string {
	res := "Project :\n"
	add := func(key, value string) {
		res += "\t" + key + " : " + value + "\n"
	}
	add("Id", strconv.Itoa(p.Id))
	add("Client", p.Client)
	add("Name", p.Name)
	add("Risk", strconv.Itoa(p.Risk))
	add("Lead PS", p.LeadPS)
	add("Lead Dev", p.LeadDev)
	add("Status", p.Status)
	add("Type", p.Type)
	add("Forecast WorkLoad", strconv.FormatFloat(p.ForecastWL, 'f', 1, 64))
	add("Current WorkLoad", strconv.FormatFloat(p.CurrentWL, 'f', 1, 64))
	add("Comment", p.Comment)
	add("Situation", p.Situation.String())
	res += ".\n"
	return res
}

// Update updates current Projet p with all element from given project p2 (p.Id will not change)
func (p *Project) Update(p2 *Project) {
	p.Client = p2.Client
	p.Name = p2.Name
	p.Risk = p2.Risk
	p.LeadPS = p2.LeadPS
	p.LeadDev = p2.LeadDev
	p.Status = p2.Status
	p.Type = p2.Type
	p.ForecastWL = p2.ForecastWL
	p.CurrentWL = p2.CurrentWL
	p.Comment = p2.Comment
	p.Situation.Update(p2.Situation.GetSituationToDate())
}
