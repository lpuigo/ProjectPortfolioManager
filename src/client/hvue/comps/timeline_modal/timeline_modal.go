package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/prjptf/src/client/business"
	fm "github.com/lpuig/prjptf/src/client/frontmodel"
	"github.com/lpuig/prjptf/src/client/tools"
	"github.com/lpuig/prjptf/src/client/tools/dates"
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
	SlotType   string  `js:"slotType"`
}

func NewTimeLineModalModel(vm *hvue.VM) *TimeLineModalModel {
	tlmm := &TimeLineModalModel{Object: tools.O()}
	tlmm.VM = vm
	tlmm.Visible = false
	tlmm.Projects = nil
	tlmm.TimeLines = nil
	tlmm.BeginDate = ""
	tlmm.EndDate = ""
	tlmm.SlotLength = 0
	tlmm.SlotType = "2Q"

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
	//tlmm.SetTimePeriod("2018-01-01", "2018-12-31")
	tlmm.HandleSlotType(tlmm.SlotType)

	tlmm.Visible = true
}

func (tlmm *TimeLineModalModel) Hide() {
	tlmm.Visible = false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SlotType methods

func (tlmm *TimeLineModalModel) HandleSlotType(slotType string) {
	tlmm.SlotType = slotType
	tlmm.setSlotDates()
	tlmm.CalcTimeLines(func(p *fm.Project) bool {
		return p.Name == "Run"
	})
}

func (tlmm *TimeLineModalModel) setSlotDates() {
	var halfPeriodLength int
	switch tlmm.SlotType {
	case "4Q":
		halfPeriodLength = 182 // half a year
	case "2Q":
		halfPeriodLength = 91 // halt a semester
	default:
		halfPeriodLength = 46 // half a quarter
	}

	tlmm.SetTimePeriod(
		date.TodayAfter(-halfPeriodLength),
		date.TodayAfter(halfPeriodLength+1),
	)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (tlmm *TimeLineModalModel) SelectRow(vm *hvue.VM, t *TimeLine, event *js.Object) {
	vm.Call("Hide")
	vm.Emit("edit-project", t.Project)
}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TimeLine Methods

func (tlmm *TimeLineModalModel) SetTimePeriod(beg, end string) {
	tlmm.BeginDate = beg
	tlmm.EndDate = end
	tlmm.SlotLength = date.NbDaysBetween(beg, end)
}

func (tlmm *TimeLineModalModel) CalcTimeLines(exclude func(*fm.Project) bool) {
	timelines := []*TimeLine{}
	for _, p := range tlmm.Projects {
		if exclude(p) {
			continue
		}
		t := tlmm.GetTimeLineFrom(p)
		if t == nil {
			continue
		}
		timelines = append(timelines, t)
	}
	tlmm.TimeLines = timelines
}

func (tlmm *TimeLineModalModel) DateOffset(d string) float64 {
	do := date.NbDaysBetween(tlmm.BeginDate, d)
	return do / tlmm.SlotLength * 100
}

func (tlmm *TimeLineModalModel) IsInSlot(beg, end string) bool {
	return end >= tlmm.BeginDate && beg <= tlmm.EndDate
}

func (tlmm *TimeLineModalModel) Length(beg, end string) float64 {
	//if beg < tlmm.BeginDate {
	//	beg = tlmm.BeginDate
	//}
	return date.NbDaysBetween(beg, end) / tlmm.SlotLength * 100
}

func (tlmm *TimeLineModalModel) GetTimeLineFrom(p *fm.Project) *TimeLine {
	t := NewTimeLine(p)

	kickoffDate, kickoffFound := p.MileStones["Kickoff"]
	outlineDate, outlineFound := p.MileStones["Outline"]
	uatDate, uatFound := p.MileStones["UAT"]
	trainingDate, trainingFound := p.MileStones["Training"]
	rolloutDate, rolloutFound := p.MileStones["RollOut"]
	goliveDate, goliveFound := p.MileStones["GoLive"]
	pilotendDate, pilotendFound := p.MileStones["Pilot End"]

	lastStylePhase := func() {
		nbPhase := len(t.Phases)
		if nbPhase >= 1 {
			t.Phases[nbPhase-1].IsLast = true
		}
	}

	className := func(phaseName string) string {
		res := phaseName
		if business.InactiveProject(p.Status) {
			res += " done"
		}
		if business.LeadProject(p.Status) {
			res += " lead"
		}
		if business.MonitoredProject(p.Status) {
			res += " monitor"
		}
		return res
	}

	setComment := func(begDate, endDate string) string {
		return begDate + " to " + endDate
	}

	if !(kickoffFound || outlineFound || uatFound || trainingFound || rolloutFound || goliveFound || pilotendFound) {
		return nil
	}

	pBeg, pEnd := date.MinMax(kickoffDate, outlineDate, uatDate, trainingDate, rolloutDate, goliveDate, pilotendDate)
	if pBeg == pEnd || !tlmm.IsInSlot(pBeg, pEnd) {
		return nil
	}
	// Study phase
	firstPhase := true
	if kickoffFound && outlineFound && tlmm.IsInSlot(kickoffDate, outlineDate) && kickoffDate != outlineDate {
		ph := NewPhase(setComment(kickoffDate, outlineDate), className("study"))
		ph.SetPositionInfo(tlmm.DateOffset(kickoffDate), tlmm.Length(kickoffDate, outlineDate))
		ph.IsFirst = true
		t.AddPhase(ph)
		firstPhase = false
	}
	// Real phase
	if trainingFound && !outlineFound {
		outlineDate = trainingDate
	}
	if uatFound && !outlineFound {
		outlineDate = uatDate
	}
	if kickoffFound && !outlineFound {
		outlineDate = kickoffDate
	}
	if goliveFound && !rolloutFound {
		rolloutDate = goliveDate
	}
	if outlineDate != "" && rolloutDate != "" && tlmm.IsInSlot(outlineDate, rolloutDate) && outlineDate != rolloutDate {
		ph := NewPhase(setComment(outlineDate, rolloutDate), className("real"))
		ph.IsFirst = true
		offset := 0.0
		if firstPhase {
			offset = tlmm.DateOffset(outlineDate)
		}
		ph.SetPositionInfo(offset, tlmm.Length(outlineDate, rolloutDate))
		t.AddPhase(ph)
		firstPhase = false
	}
	// Service phase
	if rolloutFound && goliveFound && tlmm.IsInSlot(rolloutDate, goliveDate) && rolloutDate != goliveDate {
		ph := NewPhase(setComment(rolloutDate, goliveDate), className("service"))
		ph.IsFirst = firstPhase
		offset := 0.0
		if firstPhase {
			offset = tlmm.DateOffset(rolloutDate)
		}
		ph.SetPositionInfo(offset, tlmm.Length(rolloutDate, goliveDate))
		t.AddPhase(ph)
		firstPhase = false
	}
	// Pilot phase
	if rolloutFound && !goliveFound {
		goliveDate = rolloutDate
	}
	if goliveDate != "" && pilotendFound && tlmm.IsInSlot(goliveDate, pilotendDate) && goliveDate != pilotendDate {
		lastStylePhase()
		ph := NewPhase(setComment(goliveDate, pilotendDate), className("pilot"))
		ph.IsFirst = firstPhase
		offset := 0.0
		if firstPhase {
			offset = tlmm.DateOffset(goliveDate)
		}
		ph.SetPositionInfo(offset, tlmm.Length(goliveDate, pilotendDate))
		t.AddPhase(ph)
		firstPhase = false
	}
	lastStylePhase()

	offset := 0.0
	for _, ph := range t.Phases {
		offset = ph.SetStyleClass(offset)
	}

	return t
}
