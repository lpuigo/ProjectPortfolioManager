package project_stat_modal

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/hvue/tools"
)

type IssueInfo struct {
	*js.Object

	Issue     string  `js:"issue"`
	Summary   string  `js:"summary"`
	Spent     float64 `js:"spent"`
	Remaining float64 `js:"remaining"`
	Estimated float64 `js:"estimated"`
}

func NewIssueInfo(issue, summary string, s, r, e float64) *IssueInfo {
	ii := &IssueInfo{Object: tools.O()}
	ii.Issue = issue
	ii.Summary = summary
	ii.Spent = s
	ii.Remaining = r
	ii.Estimated = e
	return ii
}

func NewIssueInfoList(ps *fm.ProjectStat) []*IssueInfo {
	res := []*IssueInfo{}
	for i, issue := range ps.Issues {
		n := len(ps.TimeSpent[i]) - 1
		ii := NewIssueInfo(
			issue,
			ps.Summaries[i],
			ps.TimeSpent[i][n],
			ps.TimeRemaining[i][n],
			ps.TimeEstimated[i][n],
		)
		res = append(res, ii)
	}
	return res
}
