package jira_stat_modal

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	jsn "github.com/lpuig/novagile/src/client/frontmodel/jirastatrecord"
	"github.com/lpuig/novagile/src/client/goel/message"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/hourstree"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/projecttree"
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
		hvue.Component("project-tree", projecttree.ComponentOptions()...),
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

	ActiveTabName    string              `js:"activeTabName"`
	WeekLogsNodes    []*hourstree.Node   `js:"wlnodes"`
	ProjectLogsNodes []*projecttree.Node `js:"plnodes"`
}

func NewJiraStatModalModel(vm *hvue.VM) *JiraStatModalModel {
	jsmm := &JiraStatModalModel{Object: tools.O()}
	jsmm.Visible = false
	jsmm.VM = vm

	jsmm.ActiveTabName = ""
	jsmm.WeekLogsNodes = []*hourstree.Node{}
	jsmm.ProjectLogsNodes = []*projecttree.Node{}
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
	jsmm.WeekLogsNodes = []*hourstree.Node{}
	jsns := jsmm.callGetJiraStat("/jira/teamlogs")
	if jsns == nil {
		return
	}
	res := []*hourstree.Node{}

	var teamNode *hourstree.Node
	team := ""

	jsns.Call("forEach", func(jsn *jsn.JiraStatRecord) {
		if jsn.Team != team {
			team = jsn.Team
			teamNode = hourstree.NewNode(team, nil, 0)
			res = append(res, teamNode)
		}
		teamNode.AddChild(hourstree.NewNode(jsn.Author, jsn.HourLogs, 40))
	})

	jsmm.WeekLogsNodes = res
}

func (jsmm *JiraStatModalModel) GetProjectLogsNodes() {
	jsmm.ProjectLogsNodes = []*projecttree.Node{}
	jsns := jsmm.callGetJiraStat("/jira/projectlogs")
	if jsns == nil {
		return
	}

	res := []*projecttree.Node{}
	var teamNode, issueNode, lotClientNode *projecttree.Node
	team, issue, lotclient := "", "", "-"

	jsns.Call("forEach", func(jsn *jsn.JiraProjectLogRecord) {
		curTeam := jsn.Infos[0]
		curActor := jsn.Infos[1]
		curLotclient := jsn.Infos[2]
		curIssue := jsn.Infos[3]
		curSummary := jsn.Infos[4]

		if team != curTeam {
			if teamNode != nil {
				teamNode.Update()
			}
			team = curTeam
			lotclient = "-"
			issue = ""
			teamNode = projecttree.NewNode(curTeam)
			res = append(res, teamNode)
		}
		if lotclient != curLotclient {
			lotclient = curLotclient
			label := curLotclient[:]
			if label == "" {
				label = "<Unassigned>"
			}
			lotClientNode = projecttree.NewNode(label)
			lotClientNode.ParentRatio = true
			teamNode.AddChild(lotClientNode)
		}
		if issue != curIssue {
			issue = curIssue
			issueNode = projecttree.NewNode(curIssue)
			issueNode.SetIssueInfo(curIssue, curSummary, jsn.TotalHour)
			issueNode.ParentRatio = true
			lotClientNode.AddChild(issueNode)
		}

		aNode := projecttree.NewNode(curActor)
		aNode.ParentRatio = true
		aNode.SetHour(curLotclient, jsn.Hour)
		issueNode.AddChild(aNode)
	})
	if teamNode != nil {
		teamNode.Update()
	}

	jsmm.ProjectLogsNodes = res

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
