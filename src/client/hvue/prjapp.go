package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_table"
	"github.com/lpuig/novagile/src/client/hvue/tools"
	"honnef.co/go/js/xhr"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#app"),
		hvue.Component("project-table", project_table.ComponentOptions()...),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.GetPtf()
		}),
	)

	// TODO to remove after debug
	js.Global.Set("mpm", mpm)

}

type MainPageModel struct {
	*js.Object

	Projects      []*fm.Project `js:"projects"`
	EditedProject *fm.Project   `js:"editedProject"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.Projects = nil
	mpm.EditedProject = nil
	return mpm
}

func (m *MainPageModel) GetPtf() {
	go m.callGetPtf()
}

func (m *MainPageModel) EditProject(p *fm.Project) {
	m.EditedProject = p
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
