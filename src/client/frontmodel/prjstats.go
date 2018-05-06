package frontmodel

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/novagile/src/client/tools"
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
		if i == 0 {
			continue // Skip first Issue as it is project global
		}
		is := NewIssueStat()
		is.Issue = issue + " : " + ps.Summaries[i]
		is.HRef = tools.UrlJiraBrowseIssue + issue
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
	sis.Issue = "Project overall"
	sis.TimeSpent = ps.TimeSpent[0]
	sis.TimeRemaining = ps.TimeRemaining[0]
	sis.TimeEstimated = ps.TimeEstimated[0]
	return sis
}
