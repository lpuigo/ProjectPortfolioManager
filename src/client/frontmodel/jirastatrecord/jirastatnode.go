package jirastatrecord

import "github.com/gopherjs/gopherjs/js"

type JiraStatRecord struct {
	*js.Object

	Team     string    `json:"team"      js:"team"`
	Author   string    `json:"author"    js:"author"`
	HourLogs []float64 `json:"hour_logs" js:"hour_logs"`
}

func NewBEJiraStatRecord(t, a string, nbweek int) *JiraStatRecord {
	return &JiraStatRecord{Team: t, Author: a, HourLogs: make([]float64, nbweek)}
}

type JiraProjectLogRecord struct {
	*js.Object

	Infos     []string `json:"infos"      js:"infos"`
	TotalHour float64  `json:"hour_total" js:"hour_total"`
	Hour      float64  `json:"hour_log"   js:"hour_log"`
}

func NewBEJiraProjectLogRecord(infos []string, totalhour, hour float64) *JiraProjectLogRecord {
	return &JiraProjectLogRecord{Infos: infos, TotalHour: totalhour, Hour: hour}
}
