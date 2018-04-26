package jiramanager

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	jsr "github.com/lpuig/novagile/src/client/frontmodel/jirastatrecord"
	"github.com/lpuig/novagile/src/server/manager/jiramanager/teamlogs"
)

type JiraManager struct {
	// TODO Jira DB attribute
	db *sql.DB
}

func NewJiraManager(usrpwd, dbname string) (*JiraManager, error) {
	res := &JiraManager{}
	jiraDb, err := sql.Open(
		"mysql",
		//"UserPwd"@"DbName",
		usrpwd+"@"+dbname,
	)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %s", err.Error())
	}
	res.db = jiraDb
	return res, nil
}

func (jm *JiraManager) TeamLogs() (jsns []*jsr.JiraStatRecord, err error) {
	jsns, err = teamlogs.Request(jm.db)
	return
}
