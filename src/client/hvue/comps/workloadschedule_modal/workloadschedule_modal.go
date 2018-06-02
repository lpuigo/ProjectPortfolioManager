package workloadschedule_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/prjptf/src/client/business"
	wsr "github.com/lpuig/prjptf/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/prjptf/src/client/goel/message"
	"github.com/lpuig/prjptf/src/client/hvue/comps/workloadschedule_modal/bars_chart"
	"github.com/lpuig/prjptf/src/client/hvue/comps/workloadschedule_modal/selectiontree"
	"github.com/lpuig/prjptf/src/client/tools"
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
		hvue.Component("selection-tree", selectiontree.ComponentOptions()...),
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
	BarChartInfos   *bars_chart.Infos     `js:"barchartInfos"`
}

func NewWSModalModel(vm *hvue.VM) *WSModalModel {
	wsmm := &WSModalModel{Object: tools.O()}
	wsmm.Visible = false
	wsmm.VM = vm

	wsmm.WrkSched = nil
	wsmm.WrkSchedLoading = false
	wsmm.BarChartInfos = bars_chart.NewInfos(nil, nil)

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
		wsmm.initSelectionTree()
		wsmm.updateBarChartInfos()
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(wsmm.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Selection Tree Methods

func (wsmm *WSModalModel) initSelectionTree() {
	wsmm.VM.Refs("selection-tree").Call("Init", wsmm.WrkSched)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Bar Chart Methods

func (wsmm *WSModalModel) updateBarChartInfos() {
	wsmm.BarChartInfos.Weeks = wsmm.WrkSched.Weeks
	wsmm.BarChartInfos.Series = []*bars_chart.Serie{}
	for _, r := range wsmm.WrkSched.Records {
		if !r.Display {
			continue
		}
		color := business.GetColorFromStatus(r.Status)
		s := bars_chart.NewSerie(r.Name, color, r.WorkLoads)
		wsmm.BarChartInfos.Series = append(wsmm.BarChartInfos.Series, s)
	}
}

func (wsmm *WSModalModel) UpdateBarChart() {
	wsmm.updateBarChartInfos()
	wsmm.VM.Refs("bars-chart").Call("Refresh", wsmm.BarChartInfos)
}
