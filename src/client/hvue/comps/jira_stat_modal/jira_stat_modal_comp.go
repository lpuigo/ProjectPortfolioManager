package jira_stat_modal

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	jsn "github.com/lpuig/element/model/jirastatnode"
	"github.com/lpuig/novagile/src/client/goel/message"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/hourstree"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/node"
	"github.com/lpuig/novagile/src/client/tools"
	"honnef.co/go/js/xhr"
)

func Register() {
	hvue.NewComponent("jira-stat-modal",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Component("hours-tree", hourstree.ComponentOptions()...),
		//hvue.Component("tab-pane", tabpane.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewJiraStatModalModel(vm)
		}),
		hvue.MethodsOf(&JiraStatModalModel{}),
	}
}

type JiraStatModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	ActiveTabName    string            `js:"activeTabName"`
	WeekLogsNodes    []*node.HoursNode `js:"wlnodes"`
	ProjectLogsNodes []*node.HoursNode `js:"plnodes"`
}

func NewJiraStatModalModel(vm *hvue.VM) *JiraStatModalModel {
	jsmm := &JiraStatModalModel{Object: tools.O()}
	jsmm.Visible = false
	jsmm.VM = vm

	jsmm.ActiveTabName = ""
	jsmm.WeekLogsNodes = []*node.HoursNode{}
	jsmm.ProjectLogsNodes = []*node.HoursNode{}
	return jsmm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (jsmm *JiraStatModalModel) Show() {
	jsmm.ActiveTabName = "weeklogs"
	jsmm.ActivateWeekLogsData()
	jsmm.Visible = true
}

func (jsmm *JiraStatModalModel) Hide() {
	jsmm.ProjectLogsNodes = nil
	jsmm.WeekLogsNodes = nil
	jsmm.Visible = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Tabs Methods

func (jsmm *JiraStatModalModel) ActivateTabs(tab *js.Object) {
	if tab == nil {
		return
	}
	tabname := tab.Get("name").String()
	println("ActivateTabs", tabname)
	switch tabname {
	case "weeklogs":
		jsmm.ActivateWeekLogsData()
	case "projectlogs":
		jsmm.ActivateProjectLogsData()
	default:

	}
}

func (jsmm *JiraStatModalModel) HandleNodeClick() {

}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (jsmm *JiraStatModalModel) ActivateWeekLogsData() {
	if len(jsmm.WeekLogsNodes) > 0 {
		return
	}
	go jsmm.GetWeekLogsNodes()
}

func (jsmm *JiraStatModalModel) ActivateProjectLogsData() {
	if len(jsmm.ProjectLogsNodes) > 0 {
		return
	}
	go jsmm.GetProjectLogsNodes()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Get Client-Name List

func (jsmm *JiraStatModalModel) GetWeekLogsNodes() {
	jsmm.WeekLogsNodes = []*node.HoursNode{}
	jsns := jsmm.callGetJiraStat("/jira/teamlogs")
	if jsns == nil {
		return
	}
	res := []*node.HoursNode{}

	var teamnode *node.HoursNode
	team := ""

	jsns.Call("forEach", func(jsn *jsn.JiraStatNode) {
		if jsn.Team != team {
			team = jsn.Team
			teamnode = node.NewHoursNode(team, nil, 0)
			res = append(res, teamnode)
		}
		teamnode.AddChild(node.NewHoursNode(jsn.Author, jsn.HourLogs, 40))
	})

	jsmm.WeekLogsNodes = res
}

func (jsmm *JiraStatModalModel) GetProjectLogsNodes() {
	jsns := jsmm.callGetJiraStat("/jira/projectlogs")
	if jsns == nil {
		return
	}
}

func (jsmm *JiraStatModalModel) callGetJiraStat(route string) *js.Object {
	//req := xhr.NewRequest("GET", "/jira/teamlogs")
	req := xhr.NewRequest("GET", route)
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(jsmm.VM, "Oups! "+err.Error(), true)
		return nil
	}
	if req.Status != 200 {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(jsmm.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
		return nil
	}
	return req.Response
}
