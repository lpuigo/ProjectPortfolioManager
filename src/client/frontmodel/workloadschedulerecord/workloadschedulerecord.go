package workloadschedulerecord

import (
	"github.com/gopherjs/gopherjs/js"
)

type WorkloadScheduleRecord struct {
	*js.Object

	Id        int       `json:"id"        js:"id"`
	Name      string    `json:"name"      js:"name"`
	Status    string    `json:"status"    js:"status"`
	LeadDev   string    `json:"leadDev"   js:"leadDev"`
	LeadPS    string    `json:"leadPS"    js:"leadPS"`
	WorkLoads []float64 `json:"workloads" js:"workloads"`
}

func NewBEWorkloadScheduleRecord(id, nbweek int, name, status, dev, ps string) *WorkloadScheduleRecord {
	return &WorkloadScheduleRecord{
		Id:        id,
		Name:      name,
		Status:    status,
		LeadDev:   dev,
		LeadPS:    ps,
		WorkLoads: make([]float64, nbweek),
	}
}

type WorkloadSchedule struct {
	*js.Object

	Weeks   []string                  `json:"weeks"   js:"weeks"`
	Records []*WorkloadScheduleRecord `json:"records" js:"records"`
}

func NewBEWorkloadSchedule(weeks []string) *WorkloadSchedule {
	return &WorkloadSchedule{Weeks: weeks, Records: []*WorkloadScheduleRecord{}}
}

func NewWorkloadScheduleFromJS(o *js.Object) *WorkloadSchedule {
	return &WorkloadSchedule{Object: o}
}
