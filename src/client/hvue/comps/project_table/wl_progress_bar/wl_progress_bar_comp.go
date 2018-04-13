package wl_progress_bar

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/hvue/tools"
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
			ppbm := &ProjectProgressBarModel{Object:vm.Object}
			return ppbm.ProgressPct()
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

func (p *ProjectProgressBarModel) ProgressPct() (pct float64) {
	pp := p.Project
	if pp.ForecastWL > 0 {
		pct = pp.CurrentWL/pp.ForecastWL*100
	} else {
		// TODO Manage unknown forcast
		pct = 50
	}
	if pct > 100 {
		// TODO Manage WL Overspent
		pct = 100
	}
	return
}
