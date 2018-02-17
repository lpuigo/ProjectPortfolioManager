package frontmodel

import (
	"github.com/gopherjs/gopherjs/js"
	"strings"
)

//go:generate easyjson.exe prjstats.go

// disabled easyjson:json
type ProjectStat struct {
	*js.Object
	Issues        []string    `json:"issues"         js:"issues"`
	Summaries     []string    `json:"summaries"      js:"summaries"`
	StartDate     string      `json:"startDate"      js:"startDate"`
	TimeSpent     [][]float64 `json:"timeSpent"      js:"timeSpent"`
	TimeRemaining [][]float64 `json:"timeRemaining"  js:"timeRemaining"`
	TimeEstimated [][]float64 `json:"timeEstimated"  js:"timeEstimated"`
}

func NewProjectStat() *ProjectStat {
	ps := &ProjectStat{Object: js.Global.Get("Object").New()}
	ps.Issues = []string{}
	ps.Summaries = []string{}
	ps.StartDate = ""
	ps.TimeSpent = [][]float64{}
	ps.TimeRemaining = [][]float64{}
	ps.TimeEstimated = [][]float64{}
	return ps
}

func NewProjectStatFromJS(o *js.Object) *ProjectStat {
	ps := &ProjectStat{Object: o}
	return ps
}

type IssueStat struct {
	*js.Object
	Issue         string    `js:"issue"`
	HRef          string    `js:"href"`
	StartDate     string    `js:"startDate"`
	TimeSpent     []float64 `js:"timeSpent"`
	TimeRemaining []float64 `js:"timeRemaining"`
	TimeEstimated []float64 `js:"timeEstimated"`
}

func NewIssueStat() *IssueStat {
	is := &IssueStat{Object: js.Global.Get("Object").New()}
	is.Issue = ""
	is.HRef = ""
	is.StartDate = ""
	is.TimeSpent = []float64{}
	is.TimeRemaining = []float64{}
	is.TimeEstimated = []float64{}
	return is
}

func CreateIssueStatsFromProjectStat(ps *ProjectStat) []*IssueStat {
	res := []*IssueStat{}
	for i, issue := range ps.Issues {
		is := NewIssueStat()
		is.Issue = issue + " : " + ps.Summaries[i]
		is.HRef = "http://jira.acticall.com/browse/" + issue
		is.StartDate = ps.StartDate
		is.TimeSpent = ps.TimeSpent[i]
		is.TimeRemaining = ps.TimeRemaining[i]
		is.TimeEstimated = ps.TimeEstimated[i]
		res = append(res, is)
	}
	return res
}

func CreateSumStatFromProjectStat(ps *ProjectStat) *IssueStat {
	sis := NewIssueStat()
	sis.StartDate = ps.StartDate
	sis.Issue = "Ensemble du projet"

	for j, _ := range ps.TimeSpent[0] {
		s, r, e := 0.0, 0.0, 0.0
		for i, _ := range ps.Issues {
			s += ps.TimeSpent[i][j]
			r += ps.TimeRemaining[i][j]
			e += ps.TimeEstimated[i][j]
		}
		sis.TimeSpent = append(sis.TimeSpent, s)
		sis.TimeRemaining = append(sis.TimeRemaining, r)
		sis.TimeEstimated = append(sis.TimeEstimated, e)
	}
	return sis
}

type ProjectStatNames struct {
	*js.Object
	Clients  []string `json:"clients"   js:"clients"`
	Projects []string `json:"projects"  js:"projects"`
}

func (psn *ProjectStatNames) GetProjectStatSignatures() []*ValText {
	res := []*ValText{}
	for i, c := range psn.Clients {
		res = append(res, NewValText(c, psn.Projects[i]))
	}
	return res
}

func NewProjectStatName() *ProjectStatNames {
	psn := &ProjectStatNames{Object: js.Global.Get("Object").New()}
	psn.Clients = []string{}
	psn.Projects = []string{}
	return psn
}

func NewProjectStatNameFromJS(o *js.Object) *ProjectStatNames {
	psn := &ProjectStatNames{Object: o}
	return psn
}

func NewProjectStatNameFromList(list []string, sep string) *ProjectStatNames {
	psn := &ProjectStatNames{}
	psn.Clients = make([]string, len(list))
	psn.Projects = make([]string, len(list))
	for i, s := range list {
		elems := strings.Split(s, sep)
		if len(elems) != 2 {
			continue
		}
		psn.Clients[i] = elems[0]
		psn.Projects[i] = elems[1]
	}
	return psn
}
