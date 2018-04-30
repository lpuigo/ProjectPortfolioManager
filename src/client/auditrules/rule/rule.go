package rule

import (
	fm "github.com/lpuig/novagile/src/client/frontmodel"
)

type Rule struct {
	Audit     *fm.Audit
	AuditFunc func(project *fm.Project) bool
}

func NewRule(prio, title string, auditfunc func(project *fm.Project) bool) *Rule {
	r := &Rule{
		Audit:     fm.NewAudit(prio, title),
		AuditFunc: auditfunc,
	}
	return r
}
