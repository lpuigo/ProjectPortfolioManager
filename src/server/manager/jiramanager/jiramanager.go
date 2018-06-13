package jiramanager

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lpuig/prjptf/src/server/manager/jiramanager/projecthistorylogs"
	"github.com/lpuig/prjptf/src/server/manager/jiramanager/projectlogs"
	"github.com/lpuig/prjptf/src/server/model"

	jsr "github.com/lpuig/prjptf/src/client/frontmodel/jirastatrecord"
	"github.com/lpuig/prjptf/src/server/manager/jiramanager/teamlogs"
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

func (jm *JiraManager) ProjectLogs() (jsns []*jsr.JiraProjectLogRecord, err error) {
	jsns, err = projectlogs.Request(jm.db)
	return
}

func (jm *JiraManager) ProjectHistoryLogs(p *model.Project) (jsns []*jsr.JiraProjectLogRecord, err error) {
	jsns, err = projecthistorylogs.Request(jm.db, p)
	return
}