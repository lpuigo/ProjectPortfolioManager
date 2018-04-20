package project_stat_modal

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/hvue/tools"
)

type IssueInfo struct {
	*js.Object

	Issue      string        `js:"issue"`
	Summary    string        `js:"summary"`
	Spent      float64       `js:"spent"`
	Remaining  float64       `js:"remaining"`
	Estimated  float64       `js:"estimated"`
	ProjectPct float64       `js:"projectPct"`
	IssueStat  *fm.IssueStat `js:"issueStat"`
}

func NewIssueInfo(issue, summary string, s, r, e, prjPct float64, is *fm.IssueStat) *IssueInfo {
	ii := &IssueInfo{Object: tools.O()}
	ii.Issue = issue
	ii.Summary = summary
	ii.Spent = s
	ii.Remaining = r
	ii.Estimated = e
	ii.IssueStat = is
	ii.ProjectPct = prjPct
	return ii
}

func NewIssueInfoList(ps *fm.ProjectStat) []*IssueInfo {
	isl := fm.CreateIssueStatsFromProjectStat(ps)
	n := len(ps.TimeSpent[0]) - 1
	totSpent := 0.0
	for _, is := range isl {
		totSpent += is.TimeSpent[n]
	}
	if totSpent == 0 {
		totSpent = 0.1 // avoid div by Zero
	}
	res := []*IssueInfo{}
	for i, is := range isl {
		s := is.TimeSpent[n]
		ii := NewIssueInfo(
			ps.Issues[i],
			ps.Summaries[i],
			s,
			is.TimeRemaining[n],
			is.TimeEstimated[n],
			float64(int(s/totSpent*1000))/10,
			is,
		)
		res = append(res, ii)
	}
	return res
}
