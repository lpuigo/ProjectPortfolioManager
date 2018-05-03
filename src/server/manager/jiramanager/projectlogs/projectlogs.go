package projectlogs

import (
	"database/sql"
	"sort"
	"strconv"
	"strings"

	jsr "github.com/lpuig/novagile/src/client/frontmodel/jirastatrecord"
	ris "github.com/lpuig/novagile/src/server/manager/recordindexedset"
)

type projectLogs struct {
	Stats *ris.RecordIndexedSet
}

func newTeamLogs() *projectLogs {
	indexes := []ris.IndexDesc{}
	indexes = append(indexes, ris.NewIndexDesc("TeamAuthorWeeks", "Team", "Author", "StartWeek"))
	indexes = append(indexes, ris.NewIndexDesc("Weeks", "StartWeek"))
	res := &projectLogs{Stats: ris.NewRecordIndexedSet(indexes...)}

	return res
}

func (tl *projectLogs) loadFromFile(file string) error {
	return tl.Stats.AddCSVDataFromFile(file)
}

func (tl *projectLogs) spentHourBy(indexname string) (keys []string, values []float64, err error) {
	cs, e := tl.Stats.GetRecordColNumByName("Hours")
	if e != nil {
		return nil, nil, e
	}
	colTimeSpent := cs[0]
	keys = tl.Stats.GetIndexKeys(indexname)
	sort.Strings(keys)
	values = make([]float64, len(keys))
	var val float64
	for i, key := range keys {
		recs := tl.Stats.GetRecordsByIndexKey(indexname, key)
		val = 0.0
		for _, rec := range recs {
			if v, err := strconv.ParseFloat(rec[colTimeSpent], 64); err != nil {
				return nil, nil, err
			} else {
				val += v
			}
		}
		values[i] = val
	}
	return
}

func Request(db *sql.DB) (jsns []*jsr.JiraStatRecord, err error) {
	tl := newTeamLogs()

	tlq := newTeamLogsQuery(db)

	/*	const testFile = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Stat_Tempo\extract Tempo 2018-04-26 utf8.csv`
		err = tl.loadFromFile(testFile)
		if err != nil {
			return
		}
	*/
	err = tl.Stats.AddDataFromQuery(tlq)
	if err != nil {
		return
	}

	weeks := tl.Stats.GetIndexKeys("Weeks")
	sort.Strings(weeks)
	weekrange := map[string]int{}
	for i, w := range weeks {
		weekstr := strings.TrimLeft(w, "!")
		weekrange[weekstr] = i
	}

	nbWeeks := len(weekrange)

	keys, hours, err := tl.spentHourBy("TeamAuthorWeeks")
	if err != nil {
		return
	}

	ota := ""
	var jsn *jsr.JiraStatRecord

	for i, key := range keys {
		cols := strings.Split(key, "!")[1:]
		numweek, found := weekrange[cols[2]]
		if !found {
			continue
		}
		ta := cols[0] + "-" + cols[1]
		if ta != ota {
			jsn = jsr.NewBEJiraStatRecord(cols[0], cols[1], nbWeeks)
			jsns = append(jsns, jsn)
			ota = ta
		}
		jsn.HourLogs[numweek] = hours[i]
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Jira Query for Team Logs

const query string = `
select
	week, team, member,lot_client, issue, summary, sum(time_spent) as time_spent
from (
	SELECT
	  wkl.startdate,
		DATE_FORMAT(wkl.startdate, '%Y-%v') AS week,
		tea.name AS team,
		wkl.author AS member,
		coalesce(lc.lot_client, '') as lot_client, 
		CONCAT(prj.pkey, "-", jra.issuenum) AS issue,
		jra.summary AS summary,
		ROUND(wkl.timeworked/3600, 2) AS time_spent
	FROM
		worklog wkl
		JOIN AO_AEFED0_TEAM_MEMBER_V2 tmm ON tmm.member_key = wkl.author
		JOIN AO_AEFED0_TEAM_V2 tea ON tea.id = tmm.team_id
		JOIN jiraissue jra ON jra.id = wkl.issueid
		JOIN project prj ON prj.id = jra.project
		JOIN nodeassociation noda ON prj.id = noda.source_node_id AND noda.source_node_entity = 'Project' AND noda.sink_node_entity = 'ProjectCategory'
-- 		JOIN projectcategory prjc ON noda.sink_node_id = prjc.id AND prjc.id IN (10000, 10001, 10103, 10200, 10500, 10700, 10900) -- cname IN ('Projets Santé', 'Projets Bancaires', 'Projets Apso', 'Projets Services', 'Business Intelligence', 'NOVAGILE', 'Novagile R&D')
		LEFT JOIN (
      select
        cfv.ISSUE,
        concat(cfc.customvalue,' - ',cfp.customvalue) as lot_client 
      from customfieldvalue cfv
      inner join customfieldoption cfc on cfc.ID = cfv.PARENTKEY
      inner join customfieldoption cfp on cfp.ID = cfv.STRINGVALUE
      where 
        cfv.customfield = 12000 and cfv.PARENTKEY is not Null
     ) lc on lc.issue = wkl.issueid
	WHERE
		tea.id IN (25, 27, 28, 33) -- team R&D (id=26) removed
		AND DATE_FORMAT(wkl.startdate, '%Y-%v') = DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL 7 DAY) , '%Y-%v')
) r
group by week, team, member,lot_client, issue, summary
order by week, team, member,lot_client, issue
;
`

type teamLogsQuery struct {
	db *sql.DB
}

func newTeamLogsQuery(db *sql.DB) *teamLogsQuery {
	tlq := &teamLogsQuery{db: db}
	return tlq
}

func (tlq *teamLogsQuery) Header() []string {
	return []string{"Team", "Author", "StartWeek", "Issue", "Summary", "Hours"}
}

func (tlq *teamLogsQuery) Query() (rows *sql.Rows, err error) {
	rows, err = tlq.db.Query(query)
	return
}

func (tlq *teamLogsQuery) Scan(r *sql.Rows) ([]string, error) {
	var Team, Author, StartWeek, Issue, Summary, Hours string

	err := r.Scan(
		&Team,
		&Author,
		&StartWeek,
		&Issue,
		&Summary,
		&Hours,
	)
	if err != nil {
		return nil, err
	}
	return []string{
		Team,
		Author,
		StartWeek,
		Issue,
		Summary,
		Hours,
	}, nil
}
