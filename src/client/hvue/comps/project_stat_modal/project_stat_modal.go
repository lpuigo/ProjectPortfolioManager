package project_stat_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/prjptf/src/client/frontmodel"
	jsn "github.com/lpuig/prjptf/src/client/frontmodel/jirastatrecord"
	"github.com/lpuig/prjptf/src/client/goel/message"
	"github.com/lpuig/prjptf/src/client/hvue/comps/jira_stat_modal/projecttree"
	"github.com/lpuig/prjptf/src/client/hvue/comps/project_stat_modal/sre_chart"
	"github.com/lpuig/prjptf/src/client/tools"
	"github.com/lpuig/prjptf/src/client/tools/dates"
	"honnef.co/go/js/xhr"
	"strconv"
)

func Register() {
	hvue.NewComponent("project-stat-modal",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Component("sre-chart", sre_chart.ComponentOptions()...),
		hvue.Component("project-tree", projecttree.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewProjectStatModalModel(vm)
		}),
		hvue.MethodsOf(&ProjectStatModalModel{}),
	}
}

type ProjectStatModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	ActiveTabName    string              `js:"activeTabName"`
	Project          *fm.Project         `js:"project"`
	ProjectStat      *fm.ProjectStat     `js:"projectStat"`
	IssueStat        *fm.IssueStat       `js:"issueStat"`
	IssueInfoList    []*IssueInfo        `js:"issueInfoList"`
	ProjectLogsNodes []*projecttree.Node `js:"plnodes"`
}

func NewProjectStatModalModel(vm *hvue.VM) *ProjectStatModalModel {
	psmm := &ProjectStatModalModel{Object: tools.O()}
	psmm.Visible = false
	psmm.VM = vm

	psmm.ActiveTabName = ""
	psmm.Project = fm.NewProject()
	psmm.ProjectStat = nil
	psmm.IssueStat = nil
	psmm.IssueInfoList = nil
	psmm.ProjectLogsNodes = []*projecttree.Node{}

	return psmm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Show Hide Methods

func (psmm *ProjectStatModalModel) Show(p *fm.Project) {
	psmm.ActiveTabName = "issuelist"
	psmm.Project = p
	psmm.ProjectStat = nil
	psmm.IssueStat = nil
	psmm.IssueInfoList = nil
	go psmm.callGetProjectStat()
	psmm.Visible = true
}

func (psmm *ProjectStatModalModel) Hide() {
	psmm.Visible = false
	psmm.Project = nil
	psmm.ProjectStat = nil
	psmm.IssueStat = nil
	psmm.IssueInfoList = nil
	psmm.ProjectLogsNodes = []*projecttree.Node{}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Tabs Methods

func (psmm *ProjectStatModalModel) ActivateTabs(tab *js.Object) {
	if tab == nil {
		return
	}
	tabname := tab.Get("name").String()
	switch tabname {
	case "projectlogs":
		psmm.ActivateProjectLogsData()
	default:
	}
}

func (psmm *ProjectStatModalModel) HandleNodeClick() {

}

func (psmm *ProjectStatModalModel) ActivateProjectLogsData() {
	if len(psmm.ProjectLogsNodes) > 0 {
		return
	}
	go psmm.GetProjectLogsNodes()
}

func (psmm *ProjectStatModalModel) GetProjectLogsNodes() {
	psmm.ProjectLogsNodes = []*projecttree.Node{}
	jsns := psmm.callGetJiraStat("/jira/projectlogs/" + strconv.Itoa(psmm.Project.Id))
	if jsns == nil {
		return
	}

	lotclientNode := projecttree.NewNode(psmm.Project.Client + " - " + psmm.Project.Name)
	res := []*projecttree.Node{lotclientNode}
	var teamNode, actorNode, issueNode *projecttree.Node
	team, actor :=  "",""

	jsns.Call("forEach", func(jsn *jsn.JiraProjectLogRecord) {
		curTeam := jsn.Infos[0]
		curActor := jsn.Infos[1]
		curIssue := jsn.Infos[2]
		curSummary := jsn.Infos[3]

		if team != curTeam {
			team = curTeam
			actor = ""
			teamNode = projecttree.NewNode(curTeam)
			teamNode.ParentRatio = true
			lotclientNode.AddChild(teamNode)
		}
		if actor != curActor {
			actor = curActor
			actorNode = projecttree.NewNode(curActor[:])
			actorNode.ParentRatio = true
			teamNode.AddChild(actorNode)
		}
		issueNode = projecttree.NewNode(curIssue)
		issueNode.SetIssueInfo(curIssue, curSummary, jsn.TotalHour)
		issueNode.ParentRatio = true
		issueNode.SetHour(curIssue, jsn.Hour)
		actorNode.AddChild(issueNode)
	})
	lotclientNode.Update()

	psmm.ProjectLogsNodes = res

}

//////////////////////////////////////////////////////////////////////////////////////////////
// XHR Methods

func (psmm *ProjectStatModalModel) callGetProjectStat() {
	req := xhr.NewRequest("GET", "/stat/"+strconv.Itoa(psmm.Project.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(psmm.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		psmm.ProjectStat = fm.NewProjectStatFromJS(req.Response)
		psmm.IssueInfoList = NewIssueInfoList(psmm.ProjectStat)
		psmm.IssueStat = fm.CreateSumStatFromProjectStat(psmm.ProjectStat)
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(psmm.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (psmm *ProjectStatModalModel) callGetJiraStat(route string) *js.Object {
	//req := xhr.NewRequest("GET", "/jira/teamlogs")
	req := xhr.NewRequest("GET", route)
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(psmm.VM, "Oups! "+err.Error(), true)
		return nil
	}
	if req.Status != 200 {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(psmm.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
		return nil
	}
	return req.Response
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Others

func (psmm *ProjectStatModalModel) FormatFloat(r, c *js.Object, v float64) string {
	return date.FormatHour(v)
}
