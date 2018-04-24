package jira_stat_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/node"
	"github.com/lpuig/novagile/src/client/tools"
)

func Register() {
	hvue.NewComponent("jira-stat-modal",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("jira-stat"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewJiraStatModalModel(vm)
		}),
		hvue.MethodsOf(&JiraStatModalModel{}),
	}
}

type JiraStatModalModel struct {
	*js.Object

	Visible           bool `js:"visible"`
	VM *hvue.VM `js:"VM"`

	Nodes     []*node.HoursNode `js:"nodes"`

}

func NewJiraStatModalModel(vm *hvue.VM) *JiraStatModalModel {
	pemm := &JiraStatModalModel{Object: tools.O()}
	pemm.Visible = false
	pemm.VM = vm
	return pemm
}

func (pemm *JiraStatModalModel) Show() {
	pemm.Visible = true
}

func (pemm *JiraStatModalModel) Hide() {
	pemm.Visible = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Button Methods

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

//////////////////////////////////////////////////////////////////////////////////////////////
// Get Client-Name List

