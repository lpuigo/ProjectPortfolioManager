package auditrules

import (
	"github.com/lpuig/novagile/src/client/auditrules/rule"
	"github.com/lpuig/novagile/src/client/business"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/server/model"
)

type Auditer struct {
	Rules []*rule.Rule
}

func NewAuditer() *Auditer {
	a := &Auditer{Rules: []*rule.Rule{}}
	return a
}

func before(d, today string) bool {
	if d == "" {
		return false
	}
	if d < today {
		return true
	}
	return false
}

func after(d, today string) bool {
	if d == "" {
		return true
	}
	if d > today {
		return true
	}
	return false
}

func max(d1, d2 string) string {
	if d1 > d2 {
		return d1
	}
	return d2
}

func ongoingPrjWithoutDeliveryDate(p *fm.Project) bool {
	if !business.OnGoingProject(p.Status) {
		return false
	}
	if p.MileStones["RollOut"] == "" && p.MileStones["GoLive"] == "" {
		return true
	}
	return false
}

func ongoingPrjWithPastDeliveryDate(p *fm.Project) bool {
	if !business.OnGoingProject(p.Status) {
		return false
	}
	today := model.Today().StringJS()
	if !ongoingPrjWithoutDeliveryDate(p) && before(max(p.MileStones["RollOut"], p.MileStones["GoLive"]), today) {
		return true
	}
	return false

}

func (a *Auditer) AddAuditRules() *Auditer {
	a.Rules = append(a.Rules, rule.NewRule("1", "Ongoing Project with undefined RollOut or GoLive date", ongoingPrjWithoutDeliveryDate))
	a.Rules = append(a.Rules, rule.NewRule("1", "Ongoing Project without estimated workload", func(p *fm.Project) bool {
		if !business.OnGoingProject(p.Status) {
			return false
		}
		if p.ForecastWL == 0 {
			return true
		}
		return false
	}))
	a.Rules = append(a.Rules, rule.NewRule("2", "Ongoing Project with past RollOut / GoLive date", ongoingPrjWithPastDeliveryDate))
	a.Rules = append(a.Rules, rule.NewRule("2", "Monitored Project more than 2 weeks after RollOut / GoLive", func(p *fm.Project) bool {
		if !business.MonitoredProject(p.Status) {
			return false
		}
		if before(max(p.MileStones["RollOut"], p.MileStones["GoLive"]), model.Today().AddDays(-14).StringJS()) {
			return true
		}
		return false
	}))
	a.Rules = append(a.Rules, rule.NewRule("2", "Inactive Project still have risk declared", func(p *fm.Project) bool {
		if !business.InactiveProject(p.Status) {
			return false
		}
		if p.Risk != "0" {
			return true
		}
		return false
	}))

	return a
}

func (a *Auditer) Audit(p *fm.Project) []*fm.Audit {
	res := []*fm.Audit{}
	for _, r := range a.Rules {
		if r.AuditFunc(p) {
			res = append(res, r.Audit)
		}
	}
	return res
}
