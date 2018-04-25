package jiramanager

import (
	jsr "github.com/lpuig/novagile/src/client/frontmodel/jirastatrecord"
	"github.com/lpuig/novagile/src/server/manager/jiramanager/teamlogs"
)

type JiraManager struct {
	// TODO Jira DB attribute

}

func NewJiraManager() (*JiraManager, error) {
	res := &JiraManager{}

	return res, nil
}

func (jm *JiraManager) TeamLogs() (jsns []*jsr.JiraStatRecord, err error) {
	jsns, err = teamlogs.Request()
	return
}
