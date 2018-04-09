package manager

import (
	"github.com/lpuig/novagile/src/server/model"
	"sort"
	"strings"
)

type Project = model.Project

type lessFunc func(p1, p2 *Project) int

type multiSorter struct {
	prjs []*Project
	less []lessFunc
}

func (ms *multiSorter) Len() int {
	return len(ms.prjs)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.prjs[i], ms.prjs[j] = ms.prjs[j], ms.prjs[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := ms.prjs[i], ms.prjs[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch less(p, q) {
		case -1:
			// p < q, so we have a decision.
			return true
		case 1:
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	switch ms.less[k](p, q) {
	case -1:
		return true
	case 0:
		return true
	}
	return false
}

func (ms *multiSorter) Sort(prjList []*Project) {
	ms.prjs = prjList
	if len(ms.less) == 0 {
		return
	}
	sort.Sort(ms)
}

func ordered(less ...lessFunc) *multiSorter {
	return &multiSorter{less: less}
}

func by(field string, asc bool) lessFunc {
	cmpStr := func(s1, s2 string) int {
		v1, v2 := strings.ToLower(s1), strings.ToLower(s2)
		if v1 == v2 {
			return 0
		} else if v1 < v2 {
			return -1
		}
		return 1
	}

	cmpFlt := func(v1, v2 float64) int {
		if v1 == v2 {
			return 0
		} else if v1 < v2 {
			return -1
		}
		return 1
	}

	reverse := func(r int) int {
		if !asc {
			return -r
		}
		return r
	}

	switch field {
	case "Client":
		return func(p1, p2 *Project) int {
			return reverse(cmpStr(p1.Client, p2.Client))
		}
	case "Projet":
		return func(p1, p2 *Project) int {
			return reverse(cmpStr(p1.Name, p2.Name))
		}
	case "Développeur":
		return func(p1, p2 *Project) int {
			return reverse(cmpStr(p1.LeadDev, p2.LeadDev))
		}
	case "Statut":
		return func(p1, p2 *Project) int {
			return reverse(cmpStr(p1.Status, p2.Status))
		}
	case "Type":
		return func(p1, p2 *Project) int {
			return reverse(cmpStr(p1.Type, p2.Type))
		}
	case "Charge":
		return func(p1, p2 *Project) int {
			return reverse(cmpFlt(p1.ForecastWL, p2.ForecastWL))
		}
	case "current_wl":
		return func(p1, p2 *Project) int {
			return reverse(cmpFlt(p1.CurrentWL, p2.CurrentWL))
		}
	case "Information":
		return func(p1, p2 *Project) int {
			return reverse(cmpStr(p1.Comment, p2.Comment))
		}
		//default:
		//	return func(p1, p2 *Project) int {
		//		return reverse(cmpStr(p1.MileStones[field], p2.MileStones[field]))
		//	}
	}
	// Return neutral comparison func otherwise
	return func(p1, p2 *Project) int { return 0 }
}
