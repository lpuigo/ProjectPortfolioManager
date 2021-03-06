package projectlogs

import (
	"database/sql"
	"fmt"

	jsr "github.com/lpuig/prjptf/src/client/frontmodel/jirastatrecord"
)

func Request(db *sql.DB) (jplrs []*jsr.JiraProjectLogRecord, err error) {
	q := newQuery(db)

	res := []*jsr.JiraProjectLogRecord{}

	qrows, e := q.Query()
	if e != nil {
		err = fmt.Errorf("could not exec query: %s", e.Error())
		return
	}
	defer qrows.Close()

	var numline int = 0
	for qrows.Next() {
		info, th, h, e := q.Scan(qrows)
		if e != nil {
			err = fmt.Errorf("could not scan line %d: %s", numline, e.Error())
			return
		}

		res = append(res, jsr.NewBEJiraProjectLogRecord(info, th, h))
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
  Team, Author, lot_client, Issue, Summary, TotalHours,
  sum(Hours) as Hours
from (
  select
    t.NAME as Team,
    wl.AUTHOR as Author,
--     date_format(wl.STARTDATE, "%Y-%v") as StartWeek,
    date(wl.STARTDATE) as StartDay,
    coalesce(lc.lot_client, '') as lot_client,
    concat(p.pkey,"-", ji.issuenum) as Issue,
    ji.SUMMARY as Summary,
    wl.timeworked / 3600 as Hours,
    ji.TIMESPENT / 3600 as TotalHours
  from worklog wl
  inner join AO_AEFED0_TEAM_MEMBER_V2 tm on tm.MEMBER_KEY = wl.AUTHOR
  inner join AO_AEFED0_TEAM_V2 t on t.ID = tm.TEAM_ID
  inner join jiraissue ji on ji.ID = wl.issueid
  inner join project p on p.ID = ji.PROJECT
  LEFT JOIN (
    select
      cfv.ISSUE,
      concat(cfc.customvalue,' - ',cfp.customvalue) as lot_client 
    from customfieldvalue cfv
    inner join customfieldoption cfc on cfc.ID = cfv.PARENTKEY
    inner join customfieldoption cfp on cfp.ID = cfv.STRINGVALUE
    where 
      cfv.customfield = 12000 and cfv.PARENTKEY is not Null
  ) lc on lc.issue = wl.issueid
  where 
  t.ID in (25, 26, 27, 28, 33)
  and date_format(wl.STARTDATE, "%Y-%v") >= date_format(DATE_SUB(CURDATE(), INTERVAL 7 DAY), "%Y-%v")
  and date_format(wl.STARTDATE, "%Y-%v") <= date_format(CURDATE(), "%Y-%v")
) tmp
group by Team, lot_client, Issue, Summary, Author
order by Team, lot_client, Issue, Summary, Author, Hours desc
;
`

type query struct {
	db *sql.DB
}

func newQuery(db *sql.DB) *query {
	tlq := &query{db: db}
	return tlq
}

func (q *query) Header() []string {
	return []string{"Team", "Author", "StartWeek", "LotClient", "Issue", "Summary", "TotalHours", "Hours"}
}

func (q *query) Query() (rows *sql.Rows, err error) {
	rows, err = q.db.Query(sqlQuery)
	return
}

func (q *query) Scan(r *sql.Rows) (infos []string, totalHour float64, logHour float64, err error) {
	var Team, Author, LotClient, Issue, Summary string
	var Hour, TotalHour float64

	err = r.Scan(
		&Team,
		&Author,
		&LotClient,
		&Issue,
		&Summary,
		&TotalHour,
		&Hour,
	)
	if err != nil {
		return nil, 0, 0, err
	}
	return []string{
		Team,
		Author,
		LotClient,
		Issue,
		Summary,
	}, TotalHour, Hour, nil
}
