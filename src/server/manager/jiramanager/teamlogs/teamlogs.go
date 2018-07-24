package teamlogs

import (
	"database/sql"
	"sort"
	"strconv"
	"strings"

	jsr "github.com/lpuig/prjptf/src/client/frontmodel/jirastatrecord"
	ris "github.com/lpuig/prjptf/src/server/manager/recordindexedset"
)

type teamLogs struct {
	Stats *ris.RecordIndexedSet
}

func newTeamLogs() *teamLogs {
	indexes := []ris.IndexDesc{}
	indexes = append(indexes, ris.NewIndexDesc("TeamAuthorWeeks", "Team", "Author", "StartWeek"))
	indexes = append(indexes, ris.NewIndexDesc("Weeks", "StartWeek"))
	res := &teamLogs{Stats: ris.NewRecordIndexedSet(indexes...)}

	return res
}

func (tl *teamLogs) loadFromFile(file string) error {
	return tl.Stats.AddCSVDataFromFile(file)
}

func (tl *teamLogs) spentHourBy(indexname string) (keys []string, values []float64, err error) {
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
	Team, Author, StartWeek, Issue, Summary,
	sum(Hours) as Hours
from (
	select
		t.NAME as Team,
		wl.AUTHOR as Author,
		date_format(wl.STARTDATE, "%Y-%v") as StartWeek,
		date(wl.STARTDATE) as StartDay,
		concat(p.pkey,"-", ji.issuenum) as Issue,
		ji.SUMMARY as Summary,
		wl.timeworked / 3600 as Hours
	from worklog wl
	inner join AO_AEFED0_TEAM_MEMBER_V2 tm on tm.MEMBER_KEY = wl.AUTHOR
	inner join AO_AEFED0_TEAM_V2 t on t.ID = tm.TEAM_ID
	inner join jiraissue ji on ji.ID = wl.issueid
	inner join project p on p.ID = ji.PROJECT
	where 
		t.ID in (25, 26, 27, 28, 33)
    and date_format(wl.STARTDATE, "%Y-%v") >= date_format(CURDATE() - interval 182 day, "%Y-%v")  
		and date_format(wl.STARTDATE, "%Y-%v") <= date_format(CURDATE(), "%Y-%v")
) tmp
group by Team, Author, StartWeek, Issue, Summary
order by Team, Author, StartWeek, Issue, Summary, Hours desc
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
