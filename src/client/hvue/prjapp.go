package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/novagile/src/client/auditrules"
	"github.com/lpuig/novagile/src/client/business"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/goel/message"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_edit_modal"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_stat_modal"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_table"
	"github.com/lpuig/novagile/src/client/tools"
	"github.com/lpuig/novagile/src/client/tools/json"
	"honnef.co/go/js/xhr"
	"strconv"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#app"),
		hvue.Component("project-table", project_table.ComponentOptions()...),
		hvue.Component("project-edit-modal", project_edit_modal.ComponentOptions()...),
		hvue.Component("project-stat-modal", project_stat_modal.ComponentOptions()...),
		hvue.Component("jira-stat-modal", jira_stat_modal.ComponentOptions()...),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.auditer = auditrules.NewAuditer().AddAuditRules()
			mpm.GetPtf()
		}),
	)
	js.Global.Get("Vue").Call("use", "ELEMENT.lang.en")

	// TODO to remove after debug
	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	Projects      []*fm.Project `js:"projects"`
	EditedProject *fm.Project   `js:"editedProject"`
	Filter        string        `js:"filter"`
	auditer       *auditrules.Auditer

	VM *hvue.VM `js:"VM"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.Projects = nil
	mpm.EditedProject = nil
	mpm.Filter = ""
	mpm.auditer = auditrules.NewAuditer().AddAuditRules()
	return mpm
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) GetPtf() {
	go m.callGetPtf()
}

func (m *MainPageModel) EditProject(p *fm.Project) {
	m.EditedProject = p
	m.VM.Refs("ProjectEdit").Call("Show", p)
}

func (m *MainPageModel) CreateNewProject() {
	p := fm.NewProject()
	p.Status = business.DefaultStatus()
	p.Type = business.DefaultType()
	p.Risk = business.DefaultRisk()
	m.EditProject(p)
}

func (m *MainPageModel) ProcessEditedProject(p *fm.Project) {
	m.EditedProject = p
	if p.Id >= 0 {
		go m.callUpdatePrj(p)
	} else {
		go m.callCreatePrj(p)
	}
}

func (m *MainPageModel) ProcessDeleteProject(p *fm.Project) {
	m.EditedProject = p
	if m.EditedProject.Id >= 0 {
		go m.callDeletePrj(m.EditedProject)
	}
}

func (m *MainPageModel) ShowProjectStat(p *fm.Project) {
	m.EditedProject = p
	m.VM.Refs("ProjectStat").Call("Show", p)
}

func (m *MainPageModel) ShowProjectAudit(p *fm.Project) {
	infos := "Audit for " + p.Client + " - " + p.Name + ":\n"
	for _, a := range p.Audits {
		infos += a.Priority + " " + a.Title + "\n"
	}
	message.InfoStr(m.VM, infos, true)
}

func (m *MainPageModel) ShowJiraStat() {
	m.VM.Refs("JiraStat").Call("Show")
}

func (m *MainPageModel) AuditProjects() {
	for _, p := range m.Projects {
		p.SetAuditResult(m.auditer.Audit(p))
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (m *MainPageModel) callGetPtf() {
	req := xhr.NewRequest("GET", "/ptf")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	//m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != 200 {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
		return
	}
	prjs := []*fm.Project{}
	req.Response.Call("forEach", func(item *js.Object) {
		p := fm.ProjectFromJS(item)
		p.SetAuditResult(m.auditer.Audit(p))
		prjs = append(prjs, p)
	})
	m.Projects = prjs
	//m.DispPrj = true
	//js.Global.Set("resp", m.Resp)
}

func (m *MainPageModel) callUpdatePrj(uprj *fm.Project) {
	req := xhr.NewRequest("PUT", "/ptf/"+strconv.Itoa(uprj.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uprj))
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		uprj.Copy(fm.ProjectFromJS(req.Response))
		uprj.SetAuditResult(m.auditer.Audit(uprj))
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (m *MainPageModel) callCreatePrj(uprj *fm.Project) {
	req := xhr.NewRequest("POST", "/ptf")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uprj))
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 201 {
		m.EditedProject.Copy(fm.ProjectFromJS(req.Response))
		m.EditedProject.SetAuditResult(m.auditer.Audit(m.EditedProject))
		m.Projects = append(m.Projects, m.EditedProject)
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (m *MainPageModel) callDeletePrj(dprj *fm.Project) {
	req := xhr.NewRequest("DELETE", "/ptf/"+strconv.Itoa(dprj.Id))
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		m.deletePrj(dprj)
		message.SetDuration(tools.SuccessMsgDuration)
		message.SuccesStr(m.VM, "Project deleted !", true)
	} else {
		message.SetDuration(tools.WarningMsgDuration)
		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (m *MainPageModel) deletePrj(dprj *fm.Project) {
	for i, p := range m.Projects {
		if p.Id == dprj.Id {
			m.EditedProject = nil
			m.Projects = append(m.Projects[:i], m.Projects[i+1:]...)
			break
		}
	}
}
