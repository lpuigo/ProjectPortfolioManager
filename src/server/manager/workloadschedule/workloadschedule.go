package workloadschedule

import (
	"github.com/lpuig/novagile/src/client/business"
	wsr "github.com/lpuig/novagile/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/novagile/src/server/model"
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
	today := model.Today()
	lastMonday := today.GetMonday()
	currentW := lastMonday.AddDays(-nbBefore * 7)
	beginDate = currentW.StringJS()
	endDate = lastMonday.AddDays((nbAfter + 1) * 7).StringJS()

	weeks = make([]string, nbBefore+nbAfter+1)
	weeks[0] = currentW.String()
	for i := 1; i < nbBefore+nbAfter+1; i++ {
		currentW = currentW.AddDays(7)
		weeks[i] = currentW.String()
	}
	return
}

// calcProjectWorkloadSchedule returns a WSRecord if given project matches the given (beg <-> end) time span, nil otherwise
func calcProjectWorkloadSchedule(project *model.Project, beg, end string, nbweeks int) *wsr.WorkloadScheduleRecord {
	curSit := project.Situation.GetSituationToDate()
	pStart := model.MaxDate(curSit.GetDatesFromKeys(business.StartMilestoneKeys()...)...)
	if pStart.IsZero() {
		// project has no start date defined => ignore it
		return nil
	}
	pEnd := model.MaxDate(curSit.GetDatesFromKeys(business.GoLiveMilestoneKeys()...)...)
	if pEnd.IsZero() {
		// project has no end date defined => ignore it
		return nil
	}
	pStartS := pStart.StringJS()
	pEndS := pEnd.StringJS()
	if pStartS > end || pEndS < beg {
		// project starts after periods, or ends before period => ignore it
		return nil
	}

	pDuration := pEnd.OpenDaysSince(pStart) + 1
	wlFactor := project.ForecastWL / float64(pDuration)

	name := project.Client + " - " + project.Name
	res := wsr.NewBEWorkloadScheduleRecord(project.Id, nbweeks, name, project.Status, project.LeadDev, project.LeadPS)

	mondayD, _ := model.DateFromJSString(beg)
	mondayS := mondayD.StringJS()
	for i := 0; i < nbweeks; i++ {
		fridayD := mondayD.AddDays(4)
		fridayS := fridayD.StringJS()
		res.WorkLoads[i] = calcWeekCoverage(mondayS, fridayS, pStartS, pEndS) * wlFactor
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
		return 5
	}
	wBd, _ := model.DateFromJSString(wBeg)
	pBd, _ := model.DateFromJSString(pBeg)
	ndB := pBd.DaysSince(wBd)
	if wBeg < pBeg && wEnd <= pEnd {
		return float64(5 - ndB)
	}
	wEd, _ := model.DateFromJSString(wEnd)
	pEd, _ := model.DateFromJSString(pEnd)
	ndE := wEd.DaysSince(pEd)
	if wBeg >= pBeg && wEnd > pEnd {
		return float64(5 - ndE)
	}
	return float64(5 - ndB - ndE)
}
