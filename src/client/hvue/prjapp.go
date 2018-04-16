package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/goel/message"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_edit_modal"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_table"
	"github.com/lpuig/novagile/src/client/hvue/tools"
	"github.com/oskca/gopherjs-json"
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
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
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

	VM *hvue.VM `js:"VM"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.Projects = nil
	mpm.EditedProject = nil
	mpm.Filter = ""
	return mpm
}

func (m *MainPageModel) GetPtf() {
	go m.callGetPtf()
}

func (m *MainPageModel) EditProject(p *fm.Project) {
	m.EditedProject = p
	m.VM.Refs("ProjectEdit").Call("Show", p)
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

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	SuccessMsgDuration = 1000
	WarningMsgDuration = 5000
)

func (m *MainPageModel) callGetPtf() {
	req := xhr.NewRequest("GET", "/ptf")
	req.Timeout = 2000
	req.ResponseType = xhr.JSON
	//m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		println("Req went wrong : ", err, req.Status)
	}
	m.Projects = []*fm.Project{}
	req.Response.Call("forEach", func(item *js.Object) {
		p := fm.ProjectFromJS(item)
		m.Projects = append(m.Projects, p)
	})
	//m.DispPrj = true
	//js.Global.Set("resp", m.Resp)
}

func (m *MainPageModel) callUpdatePrj(uprj *fm.Project) {
	req := xhr.NewRequest("PUT", "/ptf/"+strconv.Itoa(uprj.Id))
	req.Timeout = 1000
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uprj))
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		uprj.Copy(fm.ProjectFromJS(req.Response))
		message.SetDuration(SuccessMsgDuration)
		message.SuccesStr(m.VM, "Project updated !", true)
	} else {
		message.SetDuration(WarningMsgDuration)
		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (m *MainPageModel) callCreatePrj(uprj *fm.Project) {
	req := xhr.NewRequest("POST", "/ptf")
	req.Timeout = 1000
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uprj))
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 201 {
		m.EditedProject.Copy(fm.ProjectFromJS(req.Response))
		m.Projects = append(m.Projects, m.EditedProject)
		message.SetDuration(SuccessMsgDuration)
		message.SuccesStr(m.VM, "New project added !", true)
	} else {
		message.SetDuration(WarningMsgDuration)
		message.WarningStr(m.VM, "Something went wrong!\nServer returned code "+strconv.Itoa(req.Status), true)
	}
}

func (m *MainPageModel) callDeletePrj(dprj *fm.Project) {
	req := xhr.NewRequest("DELETE", "/ptf/"+strconv.Itoa(dprj.Id))
	req.Timeout = 1000
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == 200 {
		m.deletePrj(dprj)
		message.SetDuration(SuccessMsgDuration)
		message.SuccesStr(m.VM, "Project deleted !", true)
	} else {
		message.SetDuration(WarningMsgDuration)
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