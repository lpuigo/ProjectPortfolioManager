package wl_progress_bar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/tools"
	"strconv"
)

func Register() {
	hvue.NewComponent("project-progress-bar",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("project"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewProjectProgressBarModel(vm)
		}),
		hvue.MethodsOf(&ProjectProgressBarModel{}),
		hvue.Computed("progressPct", func(vm *hvue.VM) interface{} {
			ppbm := &ProjectProgressBarModel{Object: vm.Object}
			return ppbm.ProgressPct()
		}),
		hvue.Computed("showProgressBar", func(vm *hvue.VM) interface{} {
			ppbm := &ProjectProgressBarModel{Object: vm.Object}
			return ppbm.Project.ForecastWL > 0 // && ppbm.Project.CurrentWL > 0
		}),
		hvue.Computed("progressStatus", func(vm *hvue.VM) interface{} {
			ppbm := &ProjectProgressBarModel{Object: vm.Object}
			return ppbm.ProgressStatus()
		}),
		hvue.Filter("chargeFormat", func(vm *hvue.VM, prj *js.Object, args ...*js.Object) interface{} {
			return ProjectChargeFormat(fm.ProjectFromJS(prj))
		}),
	}
}

type ProjectProgressBarModel struct {
	*js.Object

	Project  *fm.Project `js:"project"`
	Progress float64     `js:"progress"`

	VM *hvue.VM `js:"VM"`
}

func NewProjectProgressBarModel(vm *hvue.VM) *ProjectProgressBarModel {
	ptm := &ProjectProgressBarModel{Object: tools.O()}
	ptm.Project = nil
	ptm.Progress = 0
	ptm.VM = vm
	return ptm
}

// Computed related Methods
//

func (pb *ProjectProgressBarModel) ProgressPct() (pct float64) {
	pp := pb.Project
	if pp.ForecastWL > 0 {
		pct = pp.CurrentWL / pp.ForecastWL * 100
	} else {
		pct = 50
	}
	if pct > 100 {
		pct = 100
	}
	return
}
func (pb *ProjectProgressBarModel) ProgressStatus() (res string) {
	pp := pb.Project
	res = "progress-bar"
	if pp.ForecastWL > 0 && pp.CurrentWL > pp.ForecastWL {
		res += " over-spent"
	}
	return
}

// Filter related Funcs
//

func ProjectChargeFormat(p *fm.Project) (res string) {
	if p.ForecastWL > 0 || p.CurrentWL > 0 {
		res = strconv.FormatFloat(p.CurrentWL, 'f', 1, 64)
		res += " / "
		res += strconv.FormatFloat(p.ForecastWL, 'f', 1, 64)
		return
	}
	res = "-"
	return
}
