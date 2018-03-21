package jirastat

import (
	"fmt"
	ris "github.com/lpuig/novagile/manager/recordindexedset"
	"strconv"
	"strings"
)

type JiraStat struct {
	Stats *ris.RecordIndexedSet
}

func NewJiraStat() *JiraStat {
	indexes := []ris.IndexDesc{}
	indexes = append(indexes, ris.NewIndexDesc("LotClientTrackedday", "CLIENT!PROJECT", "WORK_DATE"))
	indexes = append(indexes, ris.NewIndexDesc("ActorTrackedday", "ACTOR", "WORK_DATE"))
	res := &JiraStat{Stats: ris.NewRecordIndexedSet(indexes...)}

	return res
}

func (js *JiraStat) LoadFromFile(file string) error {
	return js.Stats.AddCSVDataFromFile(file)
}

func (js *JiraStat) GroupBy(indexname string) (*ris.RecordIndexedSet, error) {
	cs, err := js.Stats.GetRecordColNumByName("TIME_SPENT")
	if err != nil {
		return nil, err
	}
	colTimeSpent := cs[0]
	res := ris.NewRecordIndexedSet()
	colnames := js.Stats.GetIndexHeader(indexname)
	res.AddHeader(append(colnames, "value"))
	keys := js.Stats.GetIndexKeys(indexname)
	var val float64
	for _, key := range keys {
		groupbyvalues := strings.Split(strings.TrimLeft(key, "!"), "!")
		recs := js.Stats.GetRecordsByIndexKey(indexname, key)
		val = 0.0
		for _, rec := range recs {
			if v, err := strconv.ParseFloat(rec[colTimeSpent], 64); err != nil {
				return nil, err
			} else {
				val += v
			}
		}
		fmt.Println(groupbyvalues, strconv.FormatFloat(val, 'f', 4, 64))
	}
	return nil, nil
}
