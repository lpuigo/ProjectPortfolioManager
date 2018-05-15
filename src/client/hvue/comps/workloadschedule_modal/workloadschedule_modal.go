package workloadschedule_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/novagile/src/client/business"
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

	WrkSched        *wsr.WorkloadSchedule `js:"wrkSched"`
	WrkSchedLoading bool                  `js:"wrkSchedLoading"`
	Series          []*bars_chart.Serie   `js:"series"`
}

func NewWSModalModel(vm *hvue.VM) *WSModalModel {
	wsmm := &WSModalModel{Object: tools.O()}
	wsmm.Visible = false
	wsmm.VM = vm

	wsmm.WrkSched = nil
	wsmm.WrkSchedLoading = false
	wsmm.Series = nil

	return wsmm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Show Hide Methods

func (wsmm *WSModalModel) Show() {
	wsmm.Visible = true
	go wsmm.callGetWorkloadSchedule()
}

func (wsmm *WSModalModel) Hide() {
	wsmm.Visible = false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Others

func (wsmm *WSModalModel) callGetWorkloadSchedule() {
	wsmm.WrkSchedLoading = true
	wsmm.WrkSched = nil
	wsmm.Series = nil
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
		wsmm.calcBarsData()
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(wsmm.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (wsmm *WSModalModel) calcBarsData() {
	wsmm.Series = []*bars_chart.Serie{}
	for _, r := range wsmm.WrkSched.Records {
		color := business.GetColorFromStatus(r.Status)
		s := bars_chart.NewSerie(r.Name, color, r.WorkLoads)
		wsmm.Series = append(wsmm.Series, s)
	}
	return
}
