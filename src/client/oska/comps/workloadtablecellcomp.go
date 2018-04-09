package comps

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/oskca/gopherjs-vue"
)

type WorkLoadCellComp struct {
	*js.Object
	Project *fm.Project `js:"project"`
	GaugeOn bool        `js:"gaugeon"`
}

func NewWorkLoadCellComp() *WorkLoadCellComp {
	wlc := &WorkLoadCellComp{Object: js.Global.Get("Object").New()}
	wlc.Project = nil
	wlc.GaugeOn = false
	return wlc
}

func (wlc *WorkLoadCellComp) UpdateProgress(vm *vue.ViewModel) {
	if !wlc.GaugeOn {
		jquery.NewJQuery(vm.El.Get("firstChild")).
			Call("progress", js.M{
				"showActivity": false,
				"precision":    2,
				"label":        "ratio",
				"text": js.M{
					//"ratio": "{value} / {total}",
					"ratio": "",
				},
			})
		wlc.GaugeOn = true
	}
	jquery.NewJQuery(vm.El.Get("firstChild")).
		Call("progress", "set total", wlc.Project.ForecastWL).
		Call("progress", "set progress", wlc.Project.CurrentWL)
}

// RegisterDateTableCellComp registers to current vue intance a DateTableCell component
// having the following profile
//  <td is="workload-cell" :project="some_project"></td>
func RegisterWorkLoadCellComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewWorkLoadCellComp

	o.AddProp("project")

	//<td class="collapsing center aligned disabled">{{prj.milestones.Kickoff | DateFormat}}</td>
	o.Template = `
	<td class="center aligned" style="padding: 0px 5px">
		<div v-if="project.forecast_wl" class="ui tiny progress" :class="iserror">
			<div class="label">{{project | ChargeFormat}}</div>
			<div class="bar">
				<div class="progress"></div>
			</div>
		</div>
		<span v-else>{{project | ChargeFormat}}</span>
	</td>`

	o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		wlc := &WorkLoadCellComp{Object: vm.Object}
		if wlc.Project.ForecastWL > 0 {
			wlc.UpdateProgress(vm)
		}
	})

	o.OnLifeCycleEvent(vue.EvtUpdated, func(vm *vue.ViewModel) {
		wlc := &WorkLoadCellComp{Object: vm.Object}
		if wlc.Project.ForecastWL > 0 {
			wlc.UpdateProgress(vm)
		} else {
			wlc.GaugeOn = false
		}
	})

	o.AddComputed("iserror", func(vm *vue.ViewModel) interface{} {
		wlc := &WorkLoadCellComp{Object: vm.Object}
		if wlc.Project.ForecastWL < wlc.Project.CurrentWL {
			return "error"
		}
		return ""
	})

	return o.NewComponent().Register("workload-cell")
}
