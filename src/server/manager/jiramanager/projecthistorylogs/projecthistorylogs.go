package projecthistorylogs

import (
	"database/sql"
	"fmt"
	jsr "github.com/lpuig/prjptf/src/client/frontmodel/jirastatrecord"
	"github.com/lpuig/prjptf/src/server/model"
)

func Request(db *sql.DB, p *model.Project) (jplrs []*jsr.JiraProjectLogRecord, err error) {
	q := newQuery(db)

	res := []*jsr.JiraProjectLogRecord{}

	qrows, e := q.Query(p)
	if e != nil {
		err = fmt.Errorf("could not exec query: %s", e.Error())
		return
	}
	defer qrows.Close()

	var numline int = 0
	for qrows.Next() {
		info, _, h, e := q.Scan(qrows)
		if e != nil {
			err = fmt.Errorf("could not scan line %d: %s", numline, e.Error())
			return
		}

		res = append(res, jsr.NewBEJiraProjectLogRecord(info, 0, h))
		numline++
	}
	err = qrows.Err()
	if err != nil {
		err = fmt.Errorf("query returns: %s", err.Error())
		return
	}
	jplrs = res
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Jira Query for Project Logs

const sqlQuery string = `
select 
	team, actor, issue, summary, sum(logspent) as spent
from (
select
 	COALESCE(team.NAME, "Others Actors") as Team,
	wl.AUTHOR as Actor, 
	concat(p.pkey,'-',ji.issuenum) as Issue, 
	ji.SUMMARY as summary,
	wl.timeworked/3600 as LogSpent
from customfieldvalue cfv
inner join customfieldoption cfc on cfc.ID = cfv.PARENTKEY
inner join customfieldoption cfp on cfp.ID = cfv.STRINGVALUE
inner join jiraissue ji on ji.ID = cfv.ISSUE
inner join project p on p.ID = ji.PROJECT
inner join worklog wl on wl.issueid = ji.ID
left outer join (
	select t.NAME, tm.MEMBER_KEY
	from AO_AEFED0_TEAM_MEMBER_V2 tm 
	inner join AO_AEFED0_TEAM_V2 t on t.ID = tm.TEAM_ID and t.ID in (25, 26, 27, 28, 33)
) team on team.MEMBER_KEY = wl.AUTHOR
where 
	cfv.customfield = 12000 and cfv.PARENTKEY is not Null
      and cfc.customvalue = ?
      and cfp.customvalue = ?
) tmp
group by team, actor, issue
order by team, actor, issue
;
`

type query struct {
	db *sql.DB
}

func newQuery(db *sql.DB) *query {
	tlq := &query{db: db}
	return tlq
}

//func (q *query) Header() []string {
//	return []string{"Team", "Author", "Issue", "Summary", "Hours"}
//}

func (q *query) Query(p *model.Project) (rows *sql.Rows, err error) {
	rows, err = q.db.Query(sqlQuery, p.Client, p.Name)
	return
}

func (q *query) Scan(r *sql.Rows) (infos []string, totalHour float64, logHour float64, err error) {
	var Team, Author, Issue, Summary string
	var Hour float64

	err = r.Scan(
		&Team,
		&Author,
		&Issue,
		&Summary,
		&Hour,
	)
	if err != nil {
		return nil, 0, 0, err
	}
	return []string{
		Team,
		Author,
		Issue,
		Summary,
	}, 0, Hour, nil
}
