package model

import (
	"fmt"
	"github.com/lpuig/novagile/model/stats"
	"sort"
)

type ProjectStat struct {
	Id        int  `json:"id"`
	StartDate Date `json:"dates"`
	*stats.Stat
}

func NewProjectStat() *ProjectStat {
	ps := &ProjectStat{
		Stat:      stats.NewStat(),
		StartDate: Today(),
	}

	return ps
}

func (ps *ProjectStat) AddValues(d Date, spent, remaining, estimated float64) {
	offset := d.DaysSince(ps.StartDate)

	ps.SetValue("Spent", offset, spent)
	ps.SetValue("Remaining", offset, remaining)
	ps.SetValue("Estimated", offset, estimated)
}

func (ps *ProjectStat) String() string {
	res := fmt.Sprintf("ProjectStat (Id: %d):\n", ps.Id)
	serials := []string{}
	for k, _ := range ps.Values {
		serials = append(serials, k)
	}
	sort.Strings(serials)
	for _, serial := range serials {
		res += "\t" + serial + ":"
		offsets := []int{}
		for offset, _ := range ps.Values[serial] {
			offsets = append(offsets, offset)
		}
		sort.Ints(offsets)
		vals := ps.Values[serial]
		for _, offset := range offsets {
			res += fmt.Sprintf(" %s:%.1f", ps.StartDate.AddDays(offset).String(), vals[offset])
		}
		res += "\n"
	}
	return res
}
