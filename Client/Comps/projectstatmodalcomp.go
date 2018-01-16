package Comps

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/oskca/gopherjs-vue"
)

const (
	TemplateProjectStatModalComp = `
        <div class="ui large modal" id="ProjectStatModalComp">
			<i class="close icon"></i>
            <div class="header">
                <!--<h3>Edition du projet : <span style="color: steelblue">{{projecttitle}}</span></h3>-->
                <h3 class="ui header">
                	<i class="area chart icon"></i>
                	<div class="content">
                		Statistiques du projet : <span style="color: teal">{{givenprj.client}} - {{givenprj.name}}</span>                	
					</div>
                </h3>
            </div>

            <!--<div class="content" v-if="project">-->
            <div class="scrolling content">
            	<issue-chart :issuestat="sumstat" v-if="sumstat"></issue-chart>
            	
            	<issue-chart v-for="istat in issuestats" :issuestat="istat"></issue-chart>
            </div>
            <!--<div class="actions">-->
				<!--<div class="ui button">-->
					<!--Fermer-->
				<!--</div>-->
            <!--</div>-->
        </div>
`
)

type ProjectStatModalComp struct {
	*js.Object
	GiventPrj   *fm.Project     `js:"givenprj"`
	ProjectStat *fm.ProjectStat `js:"projectstat"`
	IssueStats  []*fm.IssueStat `js:"issuestats"`
	SumStat     *fm.IssueStat   `js:"sumstat"`
}

func NewProjectStatModalComp() *ProjectStatModalComp {
	psm := &ProjectStatModalComp{Object: js.Global.Get("Object").New()}
	psm.GiventPrj = fm.NewProject()
	psm.ProjectStat = nil
	psm.IssueStats = nil
	psm.SumStat = nil
	return psm
}

// RegisterProjectStatModalComp registers to current vue intance a ProjectStatModal component
// having the following profile
// 		<projectstat-modal
//			:givenprj="editedprj"
//			:projectstat="editedprjstat"
//		></projectstat-modal>
func RegisterProjectStatModalComp() *vue.Component {
	var jq = jquery.NewJQuery

	o := vue.NewOption()
	o.Template = TemplateProjectStatModalComp
	o.Data = NewProjectStatModalComp

	o.AddProp("givenprj", "projectstat")
	o.AddSubComponent("issue-chart", RegisterIssueChartComp())

	o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		m := &ProjectStatModalComp{Object: vm.Object}
		// setup approve and deny callback funcs
		modalOptions := js.M{
			"observeChanges": true,
			"closable":       false,
			"detachable":     true,
			"offset":         200,
			"onDeny": func() bool {
				return true
			},
			"onHidden": func() {
				m.IssueStats = nil
				m.SumStat = nil
			},
		}
		jq(vm.El).Call("modal", modalOptions)
	})

	o.AddMethod("ShowProjectStatModal", func(vm *vue.ViewModel, args []*js.Object) {
		m := &ProjectStatModalComp{Object: vm.Object}

		project := &fm.Project{Object: args[0]}
		projectStat := &fm.ProjectStat{Object: args[1]}
		m.GiventPrj = project
		m.ProjectStat = projectStat

		jq(vm.El).Call("modal", "show")
		m.IssueStats = fm.CreateIssueStatsFromProjectStat(m.ProjectStat)
		m.SumStat = fm.CreateSumStatFromProjectStat(m.ProjectStat)

		//vm.Refs.Get("IssueChart").Call("RenderChart", m.ProjectStat)
		jq(vm.El).Call("modal", "refresh")
	})

	return o.NewComponent().Register("projectstat-modal")
}
