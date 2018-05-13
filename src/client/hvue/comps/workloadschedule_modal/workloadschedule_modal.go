package workloadschedule_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/novagile/src/client/business"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	wsr "github.com/lpuig/novagile/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/novagile/src/client/goel/message"
	"github.com/lpuig/novagile/src/client/hvue/comps/workloadschedule_modal/bars_chart"
	"github.com/lpuig/novagile/src/client/tools"
	"honnef.co/go/js/xhr"
	"strconv"
)

func Register() {
	hvue.NewComponent("workloadschedule-modal",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Component("bars-chart", bars_chart.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewWSModalModel(vm)
		}),
		hvue.MethodsOf(&WSModalModel{}),
	}
}

type WSModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	Projects        []*fm.Project         `js:"projects"`
	WrkSched        *wsr.WorkloadSchedule `js:"wrkSched"`
	WrkSchedLoading bool                  `js:"wrkSchedLoading"`
	Series          []*bars_chart.Serie   `js:"series"`
}

func NewWSModalModel(vm *hvue.VM) *WSModalModel {
	wsmm := &WSModalModel{Object: tools.O()}
	wsmm.Visible = false
	wsmm.VM = vm

	wsmm.Projects = nil
	wsmm.WrkSched = nil
	wsmm.WrkSchedLoading = false
	wsmm.Series = nil

	return wsmm
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

func (wsmm *WSModalModel) Show(info *Infos) {
	wsmm.Visible = true
	wsmm.Projects = info.Projects
	go wsmm.callGetWorkloadSchedule()
}

func (wsmm *WSModalModel) Hide() {
	wsmm.Visible = false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Others

func (wsmm *WSModalModel) callGetWorkloadSchedule() {
	wsmm.WrkSchedLoading = true
	defer func() { wsmm.WrkSchedLoading = false }()

	req := xhr.NewRequest("GET", "/ptf/workload")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(wsmm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		wsmm.WrkSched = wsr.NewWorkloadScheduleFromJS(req.Response)
		wsmm.calcSeries()
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(wsmm.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (wsmm *WSModalModel) calcSeries() {
	wsmm.Series = []*bars_chart.Serie{}
	for _, r := range wsmm.WrkSched.Records {
		println("looking for", r.Id)
		name := "Not Found"
		color := "#ff3f00"
		for _, p := range wsmm.Projects {
			print(p.Id, " ")
			if p.Id == r.Id {
				name = p.Client + " - " + p.Name
				color = business.GetColorFromStatus(p.Status)
				break
			}
		}
		s := bars_chart.NewSerie(name, color, r.WorkLoads)
		wsmm.Series = append(wsmm.Series, s)
	}
	return
}
