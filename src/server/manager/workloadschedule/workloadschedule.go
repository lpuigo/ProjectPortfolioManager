package workloadschedule

import (
	wsr "github.com/lpuig/novagile/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/novagile/src/server/model"
	"time"
)

func Calc(prjs []*model.Project) *wsr.WorkloadSchedule {
	wBefore := 6
	wAfter := 6
	nbweeks := wBefore + wAfter + 1

	beginD, endD, weeks := calcWeeks(wBefore, wAfter)

	res := wsr.NewBEWorkloadSchedule(weeks)

	for _, p := range prjs {
		wsrp := calcProjectWorkloadSchedule(p, beginD, endD, nbweeks)
		if wsrp != nil {
			res.Records = append(res.Records, wsrp)
		}
	}

	return res
}

func calcWeeks(nbBefore, nbAfter int) (beginDate, endDate string, weeks []string) {
	today := model.Date(time.Now())
	lastMonday := today.GetMonday()
	currentW := lastMonday.AddDays(-nbBefore * 7)
	beginDate = currentW.StringJS()
	endDate = lastMonday.AddDays((nbAfter + 1) * 7).StringJS()

	weeks = make([]string, nbBefore+nbAfter+1)
	weeks[0] = currentW.StringWeek()
	for i := 1; i < nbBefore+nbAfter+1; i++ {
		currentW = currentW.AddDays(7)
		weeks[i] = currentW.StringWeek()
	}
	return
}

// calcProjectWorkloadSchedule returns a WSRecord if given project matches the given (beg <-> end) time span, nil otherwise
func calcProjectWorkloadSchedule(project *model.Project, beg, end string, nbweeks int) *wsr.WorkloadScheduleRecord {
	curSit := project.Situation.GetSituationToDate()
	dates := curSit.DateListJSFormat()
	if len(dates) == 0 || dates[0] > end || dates[len(dates)-1] < beg {
		// project has no date defined, or starts after periods, or ends before period => ignore it
		return nil
	}

	res := wsr.NewBEWorkloadScheduleRecord(project.Id, nbweeks)

	mondayD, _ := model.DateFromJSString(beg)
	mondayS := mondayD.StringJS()
	for i := 0; i < nbweeks; i++ {
		sundayD := mondayD.AddDays(6)
		sundayS := sundayD.StringJS()
		res.WorkLoads[i] = calcWeekCoverage(mondayS, sundayS, dates[0], dates[len(dates)-1])
		mondayD = mondayD.AddDays(7)
		mondayS = mondayD.StringJS()
	}
	return res
}

func calcWeekCoverage(wBeg, wEnd, pBeg, pEnd string) float64 {
	if wBeg > pEnd || wEnd < pBeg {
		return 0
	}
	if wBeg >= pBeg && wEnd <= pEnd {
		return 7
	}
	wBd, _ := model.DateFromJSString(wBeg)
	pBd, _ := model.DateFromJSString(pBeg)
	ndB := pBd.DaysSince(wBd)
	if wBeg < pBeg && wEnd <= pEnd {
		return float64(7 - ndB)
	}
	wEd, _ := model.DateFromJSString(wEnd)
	pEd, _ := model.DateFromJSString(pEnd)
	ndE := wEd.DaysSince(pEd)
	if wBeg >= pBeg && wEnd > pEnd {
		return float64(7 - ndE)
	}
	return float64(7 - ndB - ndE)
}
