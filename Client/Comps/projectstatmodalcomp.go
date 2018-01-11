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
                		Statistiques du projet : <span style="color: steelblue">{{givenprj.client}} - {{givenprj.name}}</span>                	
					</div>
                </h3>
            </div>

            <!--<div class="content" v-if="project">-->
            <div class="scrolling content">
            	<pre>{{projectstat}}</pre>
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
}

func NewProjectStatModalComp() *ProjectStatModalComp {
	a := &ProjectStatModalComp{Object: js.Global.Get("Object").New()}
	a.GiventPrj = fm.NewProject()
	a.ProjectStat = fm.NewProjectStat()
	return a
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
	//o.AddSubComponent("dropdown-list", RegisterDropDownListComp())

	o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		// setup approve and deny callback funcs
		modalOptions := js.M{
			"observeChanges": true,
			"closable":       false,
			"detachable":     true,
			"offset":         200,
			"onDeny": func() bool {
				return true
			},
		}
		jq(vm.El).Call("modal", modalOptions)
	})

	o.AddMethod("ShowProjectStatModal", func(vm *vue.ViewModel, args []*js.Object) {
		//m := &ProjectStatModalComp{Object: vm.Object}
		//p := &fm.Project{Object: args[0]}
		//m.GiventPrj = p
		jq(vm.El).Call("modal", "show")
		jq(vm.El).Call("modal", "refresh")
	})

	return o.NewComponent().Register("projectstat-modal")
}
