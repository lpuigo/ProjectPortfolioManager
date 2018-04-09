package comps

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/oskca/gopherjs-vue"
)

const (
	TemplateEditProjectModalComp string = `
<div class="ui modal" id="EditProjectModalComp">
    <!--<i class="close icon"></i>-->
    <div class="header">
        <h3 class="ui header">
            <i class="edit icon"></i>
            <div class="content">Edit Project : <span style="color: steelblue">{{editedprj.client}} - {{editedprj.name}}</span></div>
        </h3>
    </div>
    
    <!--<div class="content" v-if="project">-->
    <div class="scrolling content">
        <form class="ui form">
            <!--<h4 class="ui dividing header">Projet</h4>-->
            <div class="field">
                <div class="two fields">
                    <div class="field">
                        <label>Client</label>
                        <input type="text" placeholder="Client Name" v-model.trim="editedprj.client">
                    </div>
                    <div class="field">
                        <label>Project Name</label>
                        <div class="ui left action input">
                            <button class="ui icon button" @click.prevent="">
                                <!--<i class="search icon"></i>-->
                                <div ref="ProjectStatLookUpDD" class="ui dropdown search icon item" style="z-index: 1010">
                                    <i class="search icon"></i>
                                    <div class="menu">
                                        <div v-for="p in prjstatlist" class="item" @click="SetClientProject(p)">{{p.value}} - {{p.text}}</div>
                                    </div>			
                                </div>
                            </button>
                            <input type="text" placeholder="Project Name" v-model.trim="editedprj.name">
                        </div>
                    </div>
                </div>
            </div>
            <div class="fields">
                <div class="five wide field">
                    <label>PS Actor</label>
                    <input type="text" v-model.trim="editedprj.lead_ps">
                </div>
                <div class="five wide field">
                    <label>Lead Dev</label>
                    <input type="text" v-model.trim="editedprj.lead_dev">
                </div>
                <div class="four wide field">
                    <label>Type</label>
                    <dropdown-list ref="TypeDD"
                        :listvalues="types"
                        defaulttext="Project type"
                        :selected.sync="editedprj.type">
                    </dropdown-list>
                </div>
                <div class="two wide field">
                    <label>Estim. WL</label>
                    <input type="number" min="0" v-model="editedprj.forecast_wl">
                </div>
            </div>
            <div class="field">
                <label>Comment</label>
                <textarea rows="3" v-model="editedprj.comment"></textarea>
            </div>
            <div class="field">
                <div class="two fields">
                    <div class="field">
                        <label>Status</label>
                        <dropdown-list ref="StatutDD"
	                        :listvalues="statuts"
                            defaulttext="Statut du projet"
                            :selected.sync="editedprj.status">
                        </dropdown-list>
                    </div>
                    <div class="field">
                        <label>Risk</label>
                        <dropdown-list ref="RiskDD"
                            :listvalues="risks"
                            defaulttext="Risk Level"
                            :selected.sync="editedprj.risk">
                        </dropdown-list>
                    </div>
                </div>
            </div>
            <div class="field">
                <table class="ui very compact celled table">
                    <thead>
                        <tr>
                            <th class="one wide center aligned">
                                <div ref="AddMilestoneDD" class="ui dropdown icon item">
                                    <i class="big link icons">
                                        <i class="calendar outline icon"></i>
                                        <i class="corner plus green icon"></i>
                                    </i>
                                    <div class="menu">
                                        <div v-for="k in unusedMilestoneKeys" :key="k" class="item" @click="AddMilestone(k)">{{k}}</div>
                                    </div>
                                </div>
                            </th>
                            <th class="three wide right aligned">Milestone</th>
                            <th>Date</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="j in usedMilestoneKeys" :key="j">
                            <td class="center aligned">
                                <i class="big link icons" @click.prevent="DeleteMilestone(j)">
                                    <i class="calendar outline icon"></i>
                                    <i class="corner remove red icon"></i>
                                </i>
                            </td>
                            <td class="right aligned">{{j}} <i class="checked calendar icon"></i></td>
                            <td>
                                <div class="ui small input">
                                    <input type="date" v-model="editedprj.milestones[j]">
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </form>
    </div>
    <div class="actions">
        <div class="ui red right labeled icon button" v-if="editedprj.id >= 0" @click.prevent="deleteProject">
            Delete
            <i class="trash outline icon"></i>
        </div>
        <div class="ui right labeled icon button" v-if="editedprj.id >= 0" @click.prevent="duplicateProject">
            Duplicate
            <i class="clone icon"></i>
        </div>
        <div class="ui black deny button">
            Cancel
        </div>
        <div class="ui positive right labeled icon button">
            Confirm
            <i class="checkmark icon"></i>
        </div>
    </div>
</div>
`
)

type EditProjectModalComp struct {
	*js.Object
	GiventPrj   *frontmodel.Project   `js:"givenprj"`
	EditedPrj   *frontmodel.Project   `js:"editedprj"`
	PrjStatList []*frontmodel.ValText `js:"prjstatlist"`
}

func NewEditProjectModalComp() *EditProjectModalComp {
	a := &EditProjectModalComp{Object: js.Global.Get("Object").New()}
	a.GiventPrj = frontmodel.NewProject()
	a.EditedPrj = frontmodel.NewProject()
	a.PrjStatList = nil

	return a
}

// RegisterEditProjectModalComp registers to current vue intance a EditProjectModal component
// having the following profile
//  <editproject-modal
// 		:statuts="some_[]*ValText"
// 		:types="some_[]*ValText"
// 		:risks="some_[]*ValText"
//		:milestonekeys="some_[]*ValText"
// 		v-model="*Project"></editproject-modal>
func RegisterEditProjectModalComp() *vue.Component {
	var jq = jquery.NewJQuery

	o := vue.NewOption()
	o.Template = TemplateEditProjectModalComp
	o.Data = NewEditProjectModalComp

	o.AddProp("givenprj", "statuts", "types", "risks", "milestonekeys")
	o.AddSubComponent("dropdown-list", RegisterDropDownListComp())

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
			"onApprove": func() bool {
				m := &EditProjectModalComp{Object: vm.Object}
				m.GiventPrj.Copy(m.EditedPrj)
				vm.Emit("update:givenprj", m.GiventPrj)
				return true
			},
		}
		jq(vm.El).Call("modal", modalOptions)

		// Prepare dropdownlist for addmilestone
		addmilestoneDDOption := js.M{
			"on":        "hover",
			"direction": "upward",
		}
		jq(vm.Refs.Get("AddMilestoneDD")).Call("dropdown", addmilestoneDDOption)

		ProjectStatLookUpDDOption := js.M{
			"on":        "hover",
			"direction": "auto",
		}
		jq(vm.Refs.Get("ProjectStatLookUpDD")).Call("dropdown", ProjectStatLookUpDDOption)
	})

	o.AddMethod("deleteProject", func(vm *vue.ViewModel, args []*js.Object) {
		m := &EditProjectModalComp{Object: vm.Object}
		//TODO Add confirmation modal
		vm.Emit("delete:givenprj", m.GiventPrj)
		jq(vm.El).Call("modal", "hide")
	})

	o.AddMethod("duplicateProject", func(vm *vue.ViewModel, args []*js.Object) {
		m := &EditProjectModalComp{Object: vm.Object}
		m.GiventPrj = frontmodel.NewProject()
		m.EditedPrj.Id = -1
		m.EditedPrj.Name += " (Copie)"
		m.EditedPrj.CurrentWL = 0.0
		m.EditedPrj.Risk = "0"
	})

	o.AddMethod("ShowEditProjectModal", func(vm *vue.ViewModel, args []*js.Object) {
		m := &EditProjectModalComp{Object: vm.Object}
		p := &frontmodel.Project{Object: args[0]}
		m.EditedPrj.Copy(p)

		m.PrjStatList = frontmodel.NewProjectStatNameFromJS(args[1]).GetProjectStatSignatures()

		vm.Refs.Get("StatutDD").Call("changeSelected", m.EditedPrj.Status)
		vm.Refs.Get("TypeDD").Call("changeSelected", m.EditedPrj.Type)
		vm.Refs.Get("RiskDD").Call("changeSelected", m.EditedPrj.Risk)
		jq(vm.El).Call("modal", "refresh")
		jq(vm.El).Call("modal", "show")
	})

	o.AddMethod("SetClientProject", func(vm *vue.ViewModel, args []*js.Object) {
		m := &EditProjectModalComp{Object: vm.Object}
		v := &frontmodel.ValText{Object: args[0]}
		jq(vm.Refs.Get("ProjectStatLookUpDD")).Call("dropdown", "hide")
		m.EditedPrj.Client = v.Value
		m.EditedPrj.Name = v.Text
	})

	o.AddMethod("DeleteMilestone", func(vm *vue.ViewModel, args []*js.Object) {
		m := &EditProjectModalComp{Object: vm.Object}
		m.EditedPrj.RemoveMileStone(args[0].String())
	})

	o.AddMethod("AddMilestone", func(vm *vue.ViewModel, args []*js.Object) {
		m := &EditProjectModalComp{Object: vm.Object}
		jq(vm.Refs.Get("AddMilestoneDD")).Call("dropdown", "hide")
		m.EditedPrj.AddMileStone(args[0].String())
	})

	o.AddComputed("unusedMilestoneKeys", func(vm *vue.ViewModel) interface{} {
		m := &EditProjectModalComp{Object: vm.Object}
		keyList := []string{}
		vm.Get("milestonekeys").Call("forEach", func(vt *js.Object) {
			k := vt.Get("value").String()
			if _, ok := m.EditedPrj.MileStones[k]; ok == false {
				keyList = append(keyList, k)
			}
		})
		return keyList
	})

	o.AddComputed("usedMilestoneKeys", func(vm *vue.ViewModel) interface{} {
		m := &EditProjectModalComp{Object: vm.Object}
		keyList := []string{}
		vm.Get("milestonekeys").Call("forEach", func(vt *js.Object) {
			k := vt.Get("value").String()
			if _, ok := m.EditedPrj.MileStones[k]; ok == true {
				keyList = append(keyList, k)
			}
		})
		return keyList
	})

	return o.NewComponent().Register("editproject-modal")
}
