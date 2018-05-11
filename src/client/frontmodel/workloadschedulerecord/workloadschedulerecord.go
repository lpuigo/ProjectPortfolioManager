package workloadschedulerecord

import "github.com/gopherjs/gopherjs/js"

type WorkloadScheduleRecord struct {
	*js.Object

	Id        int       `json:"id"        js:"id"`
	WorkLoads []float64 `json:"workloads" js:"workloads"`
}

func NewBEWorkloadScheduleRecord(id, nbweek int) *WorkloadScheduleRecord {
	return &WorkloadScheduleRecord{Id: id, WorkLoads: make([]float64, nbweek)}
}

type WorkloadSchedule struct {
	*js.Object

	Weeks   []string                  `json:"weeks"        js:"weeks"`
	Records []*WorkloadScheduleRecord `json:"records" js:"records"`
}

func NewBEWorkloadSchedule(weeks []string) *WorkloadSchedule {
	return &WorkloadSchedule{Weeks: weeks, Records: []*WorkloadScheduleRecord{}}
}
