package project_table

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/novagile/src/client/business"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/hvue/tools"
)

func Register() {
	hvue.NewComponent("project-table",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("selected_project", "projects", "filter"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewProjectTableModel(vm)
		}),
		hvue.MethodsOf(&ProjectTableModel{}),
		hvue.Computed("filteredProjects", func(vm *hvue.VM) interface{} {
			ptm := &ProjectTableModel{Object: vm.Object}
			if ptm.Filter == "" {
				return ptm.Projects
			}
			res := []*fm.Project{}
			for _, p := range ptm.Projects {
				if fm.TextFiltered(p, ptm.Filter) {
					res = append(res, p)
				}
			}
			return res
		}),
	}
}

func NewProjectTableModel(vm *hvue.VM) *ProjectTableModel {
	ptm := &ProjectTableModel{Object: tools.O()}
	ptm.Projects = nil
	ptm.SelectedProject = nil
	ptm.Filter = ""
	ptm.VM = vm
	return ptm
}

type ProjectTableModel struct {
	*js.Object

	SelectedProject *fm.Project   `js:"selected_project"`
	Projects        []*fm.Project `js:"projects"`
	Filter          string        `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func (ptm *ProjectTableModel) SelectRow(vm *hvue.VM, prj *fm.Project, event *js.Object) {
	vm.Emit("selected_project", prj)
}

func (ptm *ProjectTableModel) SetSelectedProject(np *fm.Project) {
	//ptm = &ProjectTableCompModel{Object: vm.Object}
	ptm.SelectedProject = np
	ptm.VM.Emit("update:selected_project", np)
}

func (ptm *ProjectTableModel) TableRowClassName(rowInfo *js.Object) string {
	p := &fm.Project{Object: rowInfo.Get("row")}
	var res string
	switch p.Status {
	case "6 - Done", "0 - Lost":
		res = "project-row-done"
	case "1 - Candidate", "2 - Outlining":
		res = "project-row-outline"
	default:
		res = ""
	}
	return res
}

func (ptm *ProjectTableModel) StatusList() []*fm.ValText {
	return business.CreateStatuts()
}

func (ptm *ProjectTableModel) StatusFilter(value string, p *fm.Project) bool {
	return p.Status == value
}
