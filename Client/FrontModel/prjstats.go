package FrontModel

import "github.com/gopherjs/gopherjs/js"

type ProjectStat struct {
	*js.Object
	Issues        []string    `json:"issues"         js:"issues"`
	Dates         []string    `json:"dates"          js:"dates"`
	TimeSpent     [][]float64 `json:"timeSpent"      js:"timeSpent"`
	TimeRemaining [][]float64 `json:"timeRemaining"  js:"timeRemaining"`
	TimeEstimated [][]float64 `json:"timeEstimated"  js:"timeEstimated"`
}

func NewProjectStat() *ProjectStat {
	ps := &ProjectStat{Object: js.Global.Get("Object").New()}
	ps.Issues = []string{}
	ps.Dates = []string{}
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
	Dates         []string  `js:"dates"`
	TimeSpent     []float64 `js:"timeSpent"`
	TimeRemaining []float64 `js:"timeRemaining"`
	TimeEstimated []float64 `js:"timeEstimated"`
}

func NewIssueStat() *IssueStat {
	is := &IssueStat{Object: js.Global.Get("Object").New()}
	is.Issue = ""
	is.HRef = ""
	is.Dates = []string{}
	is.TimeSpent = []float64{}
	is.TimeRemaining = []float64{}
	is.TimeEstimated = []float64{}
	return is
}

func CreateIssueStatsFromProjectStat(ps *ProjectStat) []*IssueStat {
	res := []*IssueStat{}
	for i, issue := range ps.Issues {
		is := NewIssueStat()
		is.Issue = issue
		is.HRef = "http://jira.acticall.com/browse/" + issue
		is.Dates = ps.Dates
		is.TimeSpent = ps.TimeSpent[i]
		is.TimeRemaining = ps.TimeRemaining[i]
		is.TimeEstimated = ps.TimeEstimated[i]
		res = append(res, is)
	}
	return res
}

func CreateSumStatFromProjectStat(ps *ProjectStat) *IssueStat {
	sis := NewIssueStat()
	sis.Dates = ps.Dates
	sis.Issue = "Ensemble du projet"

	for j, _ := range ps.Dates {
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
