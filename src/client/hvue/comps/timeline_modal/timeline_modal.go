package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/tools"
	"github.com/lpuig/novagile/src/client/tools/dates"
)

func Register() {
	hvue.NewComponent("timeline-modal",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		//hvue.Component("sre-chart", sre_chart.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewTimeLineModalModel(vm)
		}),
		hvue.MethodsOf(&TimeLineModalModel{}),
	}
}

type TimeLineModalModel struct {
	*js.Object

	VM      *hvue.VM `js:"VM"`
	Visible bool     `js:"visible"`

	Projects  []*fm.Project `js:"projects"`
	TimeLines []*TimeLine   `js:"timelines"`

	BeginDate  string  `js:"beginDate"`
	EndDate    string  `js:"endDate"`
	SlotLength float64 `js:"slotLength"`
}

func NewTimeLineModalModel(vm *hvue.VM) *TimeLineModalModel {
	tlmm := &TimeLineModalModel{Object: tools.O()}
	tlmm.VM = vm
	tlmm.Visible = false
	tlmm.Projects = nil
	tlmm.BeginDate = ""
	tlmm.EndDate = ""
	tlmm.SlotLength = 0

	return tlmm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Show Hide Methods

type Infos struct {
	*js.Object

	Projects []*fm.Project `js:"projects"`
}

func NewInfos(prjs []*fm.Project) *Infos {
	i := &Infos{Object: tools.O()}
	i.Projects = prjs
	return i
}

func (tlmm *TimeLineModalModel) Show(infos *Infos) {
	tlmm.Projects = infos.Projects
	tlmm.SetTimePeriod("2018-01-01", "2018-12-31")

	//go tlmm.callGetProjectStat()
	tlmm.Visible = true
}

func (tlmm *TimeLineModalModel) Hide() {
	tlmm.Visible = false
	tlmm.Projects = nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Others

func (tlmm *TimeLineModalModel) SetTimePeriod(beg, end string) {
	tlmm.BeginDate = beg
	tlmm.EndDate = end
	tlmm.SlotLength = date.NbDaysBetween(beg, end)
}

func (tlmm *TimeLineModalModel) DateOffset(d string) float64 {
	do := date.NbDaysBetween(tlmm.BeginDate, d)
	return do / tlmm.SlotLength
}

func (tlmm *TimeLineModalModel) IsInSlot(beg, end string) bool {
	return end >= tlmm.BeginDate && beg <= tlmm.EndDate
}

func (tlmm *TimeLineModalModel) Length(beg, end string) float64 {
	return date.NbDaysBetween(beg, end) / tlmm.SlotLength
}

func (tlmm *TimeLineModalModel) GetTimeLineFrom(p *fm.Project) *TimeLine {
	t := NewTimeLine(p.Client + " - " + p.Name)

	kickoffDate, kickoffFound := p.MileStones["Kickoff"]
	outlineDate, outlineFound := p.MileStones["Outline"]
	uatDate, uatFound := p.MileStones["UAT"]
	trainingDate, trainingFound := p.MileStones["Training"]
	rolloutDate, rolloutFound := p.MileStones["RollOut"]
	goliveDate, goliveFound := p.MileStones["GoLive"]
	pilotendDate, pilotendFound := p.MileStones["Pilot End"]

	if !(kickoffFound || outlineFound || uatFound || trainingFound || rolloutFound || goliveFound || pilotendFound) {
		return nil
	}

	pBeg, pEnd := date.MinMax(kickoffDate, outlineDate, uatDate, trainingDate, rolloutDate, goliveDate, pilotendDate)
	if !tlmm.IsInSlot(pBeg, pEnd) {
		return nil
	}
	firstPhase := true
	if kickoffFound && outlineFound && tlmm.IsInSlot(kickoffDate, outlineDate) {
		p := NewPhase("study")
		p.SetStyle(tlmm.DateOffset(kickoffDate), tlmm.Length(kickoffDate, outlineDate))
		t.AddPhase(p)
		firstPhase = false
	}
	if kickoffFound && !outlineFound {
		outlineDate = kickoffDate
	}
	if goliveFound && !rolloutFound {
		rolloutDate = goliveDate
	}
	if outlineDate != "" && rolloutDate != "" && tlmm.IsInSlot(outlineDate, rolloutDate) {
		p := NewPhase("real")
		offset := 0.0
		if firstPhase {
			offset = tlmm.DateOffset(outlineDate)
		}
		p.SetStyle(offset, tlmm.Length(outlineDate, rolloutDate))
		t.AddPhase(p)
		firstPhase = false
	}
	if rolloutFound && goliveFound && tlmm.IsInSlot(rolloutDate, goliveDate) {
		p := NewPhase("service")
		offset := 0.0
		if firstPhase {
			offset = tlmm.DateOffset(rolloutDate)
		}
		p.SetStyle(offset, tlmm.Length(rolloutDate, goliveDate))
		t.AddPhase(p)
		firstPhase = false
	}
	if rolloutFound && !goliveFound {
		goliveDate = rolloutDate
	}
	if goliveDate != "" && pilotendFound && tlmm.IsInSlot(goliveDate, pilotendDate) {
		p := NewPhase("real")
		offset := 0.0
		if firstPhase {
			offset = tlmm.DateOffset(goliveDate)
		}
		p.SetStyle(offset, tlmm.Length(goliveDate, pilotendDate))
		t.AddPhase(p)
		firstPhase = false
	}

	return t
}
