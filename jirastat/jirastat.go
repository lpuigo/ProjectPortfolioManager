package jirastat

import (
	ris "github.com/lpuig/novagile/manager/recordindexedset"
	"sort"
	"strconv"
)

type JiraStat struct {
	Stats *ris.RecordIndexedSet
}

func NewJiraStat() *JiraStat {
	indexes := []ris.IndexDesc{}
	indexes = append(indexes, ris.NewIndexDesc("LotClient", "CLIENT!PROJECT"))
	indexes = append(indexes, ris.NewIndexDesc("LotClientTrackedday", "CLIENT!PROJECT", "WORK_DATE"))
	indexes = append(indexes, ris.NewIndexDesc("ActorTrackedday", "ACTOR", "WORK_DATE"))
	res := &JiraStat{Stats: ris.NewRecordIndexedSet(indexes...)}

	return res
}

func (js *JiraStat) LoadFromFile(file string) error {
	return js.Stats.AddCSVDataFromFile(file)
}

func (js *JiraStat) SpentHourBy(indexname string) (keys []string, values []float64, err error) {
	cs, e := js.Stats.GetRecordColNumByName("TIME_SPENT")
	if e != nil {
		return nil, nil, e
	}
	colTimeSpent := cs[0]
	keys = js.Stats.GetIndexKeys(indexname)
	sort.Strings(keys)
	values = make([]float64, len(keys))
	var val float64
	for i, key := range keys {
		recs := js.Stats.GetRecordsByIndexKey(indexname, key)
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
