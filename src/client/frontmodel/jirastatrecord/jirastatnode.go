package jirastatrecord

import "github.com/gopherjs/gopherjs/js"

type JiraStatRecord struct {
	*js.Object

	Team     string    `json:"team"      js:"team"`
	Author   string    `json:"author"    js:"author"`
	HourLogs []float64 `json:"hour_logs" js:"hour_logs"`
}

func NewBEJiraStatRecord(t, a string, nbweek int) *JiraStatRecord {
	return &JiraStatRecord{Team:t, Author:a, HourLogs:make([]float64,nbweek)}
}


