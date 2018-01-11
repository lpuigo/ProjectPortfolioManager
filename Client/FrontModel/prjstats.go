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
