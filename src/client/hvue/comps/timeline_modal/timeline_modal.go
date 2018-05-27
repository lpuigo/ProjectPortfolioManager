package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/tools"
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
}

func NewTimeLineModalModel(vm *hvue.VM) *TimeLineModalModel {
	tlmm := &TimeLineModalModel{Object: tools.O()}
	tlmm.VM = vm
	tlmm.Visible = false
	tlmm.Projects = nil

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

	//go tlmm.callGetProjectStat()
	tlmm.Visible = true
}

func (tlmm *TimeLineModalModel) Hide() {
	tlmm.Visible = false
	tlmm.Projects = nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Others
func (tlmm *TimeLineModalModel) GetPhases(p *fm.Project) []*TimeLine {

}
