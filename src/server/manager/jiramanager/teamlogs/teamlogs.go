package teamlogs

import (
	jsr "github.com/lpuig/novagile/src/client/frontmodel/jirastatrecord"
	ris "github.com/lpuig/novagile/src/server/manager/recordindexedset"
	"sort"
	"strconv"
	"strings"
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

func Request() (jsns []*jsr.JiraStatRecord, err error) {
	tl := newTeamLogs()
	//TODO do Query
	//panic("implement query")
	const testFile = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Stat_Tempo\extract Tempo 2018-04-26 utf8.csv`
	err = tl.loadFromFile(testFile)
	if err != nil {
		return
	}

	weeks := tl.Stats.GetIndexKeys("Weeks")
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
