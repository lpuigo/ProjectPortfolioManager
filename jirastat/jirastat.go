package jirastat

import ris "github.com/lpuig/novagile/manager/recordindexedset"

type JiraStat struct {
	*ris.RecordLinkedIndexedSet
}

func NewJiraStat() *JiraStat {
	indexes := []ris.IndexDesc{}
	indexes = append(indexes, ris.NewIndexDesc("LotClientTrackedday", "CLIENT!PROJECT", "WORK_DATE"))
	indexes = append(indexes, ris.NewIndexDesc("ActorTrackedday", "ACTOR", "WORK_DATE"))
	res := &JiraStat{RecordLinkedIndexedSet: ris.NewRecordLinkedIndexedSet(indexes...)}

	res.AddLink()

	return res
}

func (js *JiraStat) LoadFromFile(file string) error {
	return js.AddCSVDataFromFile(file)
}
