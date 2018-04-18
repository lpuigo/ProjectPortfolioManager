package project_stat_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/goel/message"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_stat_modal/sre_chart"
	"github.com/lpuig/novagile/src/client/hvue/tools"
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
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewProjectStatModalModel(vm)
		}),
		hvue.MethodsOf(&ProjectStatModalModel{}),
	}
}

type ProjectStatModalModel struct {
	*js.Object

	Visible bool `js:"visible"`

	Project       *fm.Project     `js:"project"`
	ProjectStat   *fm.ProjectStat `js:"projectStat"`
	IssueStat     *fm.IssueStat   `js:"issueStat"`
	IssueInfoList []*IssueInfo    `js:"issueInfoList"`

	VM *hvue.VM `js:"VM"`
}

func NewProjectStatModalModel(vm *hvue.VM) *ProjectStatModalModel {
	psmm := &ProjectStatModalModel{Object: tools.O()}
	psmm.Visible = false
	psmm.Project = fm.NewProject()
	psmm.ProjectStat = nil
	psmm.IssueStat = nil
	psmm.IssueInfoList = nil

	psmm.VM = vm
	return psmm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Show Hide Methods

func (psmm *ProjectStatModalModel) Show(p *fm.Project) {
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
}

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Others
