package project_table

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/hvue/comps/project_table/wl_progress_bar"
	"github.com/lpuig/novagile/src/client/hvue/tools"
	"sort"
	"strconv"
)

func Register() {
	hvue.NewComponent("project-table",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Component("project-progress-bar", wl_progress_bar.ComponentOptions()...),
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
		hvue.Filter("DateFormat", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			return fm.DateString(value.String())
		}),
	}
}

type ProjectTableModel struct {
	*js.Object

	SelectedProject *fm.Project   `js:"selected_project"`
	Projects        []*fm.Project `js:"projects"`
	Filter          string        `js:"filter"`

	VM *hvue.VM `js:"VM"`
}

func NewProjectTableModel(vm *hvue.VM) *ProjectTableModel {
	ptm := &ProjectTableModel{Object: tools.O()}
	ptm.Projects = nil
	ptm.SelectedProject = nil
	ptm.Filter = ""
	ptm.VM = vm
	return ptm
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Event Related Methods

func (ptm *ProjectTableModel) SelectRow(vm *hvue.VM, prj *fm.Project, event *js.Object) {
	vm.Emit("selected_project", prj)
}

func (ptm *ProjectTableModel) ShowTableProjectStat(vm *hvue.VM, prj *fm.Project) {
	vm.Emit("show_project_stat", prj)
}

func (ptm *ProjectTableModel) SetSelectedProject(np *fm.Project) {
	//ptm = &ProjectTableCompModel{Object: vm.Object}
	ptm.SelectedProject = np
	ptm.VM.Emit("update:selected_project", np)
}

//
// Formatting Related Methods

func (ptm *ProjectTableModel) TableRowClassName(rowInfo *js.Object) string {
	p := &fm.Project{Object: rowInfo.Get("row")}
	var res string
	switch p.Status {
	case "6 - Done", "0 - Lost":
		res = "project-row-done"
	case "5 - Monitoring":
		res = "project-row-monitoring"
	case "1 - Candidate", "2 - Outlining":
		res = "project-row-outline"
	default:
		res = ""
	}
	return res
}

func (ptm *ProjectTableModel) HeaderCellStyle() string {
	return "background: #a1e6e6;"
}

func (ptm *ProjectTableModel) RiskIconClass(risk string) string {
	var res string
	switch risk {
	case "2":
		res = "fas fa-exclamation-triangle risk-icon high-risk"
	case "1":
		res = "fas fa-exclamation-circle risk-icon low-risk"
	default:
		//res = "green info circle icon"
	}
	return res
}

func (ptm *ProjectTableModel) FormatDate(r, c *js.Object, date string) string {
	return fm.DateString(date)
}

//
// Column Filtering Related Methods

func (ptm *ProjectTableModel) FilterHandler(value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (ptm *ProjectTableModel) FilterList(vm *hvue.VM, prop string) []*fm.ValText {
	ptm = &ProjectTableModel{Object: vm.Object}
	count := map[string]int{}
	attribs := []string{}
	for _, p := range ptm.Projects {
		attrib := p.Object.Get(prop).String()
		if _, exist := count[attrib]; !exist {
			attribs = append(attribs, attrib)
		}
		count[attrib]++
	}
	sort.Strings(attribs)
	res := []*fm.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "<Empty>"
		}
		res = append(res, fm.NewValText(a, fa+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

func (ptm *ProjectTableModel) FilteredValue() []string {
	res := []string{
		"1 - Candidate",
		"2 - Outlining",
		"3 - On Going",
		"4 - UAT",
		"5 - Monitoring",
	}
	return res
}
