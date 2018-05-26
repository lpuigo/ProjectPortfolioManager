package frontmodel

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/novagile/src/client/tools"
	"github.com/lpuig/novagile/src/server/model"
	"strconv"
	"strings"
)

type Audit struct {
	*js.Object

	Title    string `js:"title"`
	Priority string `js:"priority"`
}

func NewAudit(prio, title string) *Audit {
	a := &Audit{Object: tools.O()}
	a.Priority = prio
	a.Title = title
	return a
}

type Project struct {
	*js.Object
	Id         int               `json:"id"               js:"id"`
	Client     string            `json:"client"           js:"client"`
	Name       string            `json:"name"             js:"name"`
	Risk       string            `json:"risk"             js:"risk"`
	LeadPS     string            `json:"lead_ps"          js:"lead_ps"`
	LeadDev    string            `json:"lead_dev"         js:"lead_dev"`
	Status     string            `json:"status"           js:"status"`
	Type       string            `json:"type"             js:"type"`
	HasStat    bool              `json:"hasStat"          js:"hasStat"`
	ForecastWL float64           `json:"forecast_wl"      js:"forecast_wl"`
	CurrentWL  float64           `json:"current_wl"       js:"current_wl"`
	Comment    string            `json:"comment"          js:"comment"`
	Audits     []*Audit          `json:"audits"           js:"audits"` // not used on server side
	MileStones map[string]string `json:"milestones"       js:"milestones"`
}

func NewProject() *Project {
	pf := &Project{Object: tools.O()}
	pf.Id = -1
	pf.Client = "New Client"
	pf.Name = "New Projet"
	pf.Risk = ""
	pf.LeadPS = ""
	pf.LeadDev = ""
	pf.Status = ""
	pf.Type = ""
	pf.HasStat = false
	pf.ForecastWL = 0
	pf.CurrentWL = 0
	pf.Comment = ""
	pf.MileStones = nil
	pf.Audits = nil
	return pf
}

func (p *Project) Clone() *Project {
	np := &Project{Object: js.Global.Get("Object").New()}
	np.Copy(p)
	return np
}

func (p *Project) Copy(np *Project) {
	p.Id = np.Id
	p.Client = np.Client
	p.Name = np.Name
	p.Risk = np.Risk
	p.LeadPS = np.LeadPS
	p.LeadDev = np.LeadDev
	p.Status = np.Status
	p.Type = np.Type
	p.HasStat = np.HasStat
	p.ForecastWL = np.ForecastWL
	p.CurrentWL = np.CurrentWL
	p.Comment = np.Comment
	p.Audits = np.Audits[:]

	m := make(map[string]string)
	mop := np.Get("milestones")
	for _, k := range js.Keys(mop) {
		m[k] = mop.Get(k).String()
	}
	p.MileStones = m
}

func (p Project) String() string {
	res := "FrontModel Project :\n"
	add := func(key, value string) {
		res += "\t" + key + " : " + value + "\n"
	}
	add("Id", strconv.Itoa(p.Id))
	add("Client", p.Client)
	add("Name", p.Name)
	add("Risk", p.Risk)
	add("Lead PS", p.LeadPS)
	add("Lead Dev", p.LeadDev)
	add("Status", p.Status)
	add("Type", p.Type)
	add("HasStat", strconv.FormatBool(p.HasStat))
	add("Forecast WorkLoad", strconv.FormatFloat(p.ForecastWL, 'f', 1, 64))
	add("Current WorkLoad", strconv.FormatFloat(p.CurrentWL, 'f', 1, 64))
	add("Comment", p.Comment)
	add("Situation", p.Comment)
	for k, v := range p.MileStones {
		res += "\t\t" + k + " : " + DateString(v) + "\n"
	}
	res += ".\n"
	return res
}

func DateString(v string) string {
	if strings.Contains(v, "-") {
		d := strings.Split(v, "-")
		return d[2] + "/" + d[1] + "/" + d[0]
	}
	return "-"
}

func (p *Project) SearchInString() string {
	// HasStat is skipped on purpose
	res := "Client:" + p.Client + "\n"
	res += "Name:" + p.Name + "\n"
	res += "PS:" + p.LeadPS + "\n"
	res += "Dev:" + p.LeadDev + "\n"
	res += "Status:" + p.Status + "\n"
	res += "Type:" + p.Type + "\n"
	res += p.Comment + "\n"
	for m, v := range p.MileStones {
		res += m + ":" + DateString(v) + "\n"
	}
	return res
}

func (p *Project) Contains(str string) bool {
	if str == "" {
		return true
	}
	return strings.Contains(strings.ToLower(p.SearchInString()), strings.ToLower(str))
}

func (p *Project) RemoveMileStone(msName string) {
	//p.Get("milestones").Delete(msName)
	nms := make(map[string]string)
	for k, v := range p.MileStones {
		if k == msName {
			continue
		}
		nms[k] = v
	}
	p.MileStones = nms
}

func (p *Project) AddMileStone(msName string) {
	//p.Get("milestones").Set(msName, Model.Today().StringJS())
	nms := make(map[string]string)
	for k, v := range p.MileStones {
		nms[k] = v
	}
	nms[msName] = model.Today().StringJS()
	p.MileStones = nms
}

func (p *Project) SetAuditResult(audits []*Audit) {
	p.Audits = audits[:]
}

func CloneBEProject(p *model.Project, hasStat bool) *Project {
	np := &Project{}
	np.Id = p.Id
	np.Client = p.Client
	np.Name = p.Name
	np.Risk = strconv.Itoa(p.Risk)
	np.LeadPS = p.LeadPS
	np.LeadDev = p.LeadDev
	np.Status = p.Status
	np.Type = p.Type
	np.HasStat = hasStat
	np.ForecastWL = p.ForecastWL
	np.CurrentWL = p.CurrentWL
	np.Comment = p.Comment
	np.Audits = []*Audit{} // set to empty slice in order to force JS initialisation when unmarshalling
	np.MileStones = make(map[string]string)
	for m, d := range p.Situation.GetSituationToDate().MileStones {
		np.MileStones[m] = d.StringJS()
	}
	return np
}

func CloneFEProject(p *Project) *model.Project {
	np := &model.Project{}
	np.Id = p.Id
	np.Client = p.Client
	np.Name = p.Name
	r, _ := strconv.ParseInt(p.Risk, 10, 0)
	np.Risk = int(r)
	np.LeadPS = p.LeadPS
	np.LeadDev = p.LeadDev
	np.Status = p.Status
	np.Type = p.Type
	np.ForecastWL = p.ForecastWL
	np.CurrentWL = p.CurrentWL
	np.Comment = p.Comment
	np.Situation = model.NewSituations()
	std := model.NewSituationToDate()
	for m, d := range p.MileStones {
		std.MileStones[m], _ = model.DateFromJSString(d)
	}
	np.Situation.Update(std)
	return np
}

func ProjectFromJS(o *js.Object) *Project {
	p := &Project{Object: o}
	return p
}
