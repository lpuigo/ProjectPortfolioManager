package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	biz "github.com/lpuig/novagile/client/business"
	"github.com/lpuig/novagile/client/oska/comps"
	fm "github.com/lpuig/novagile/client/frontmodel"
	"github.com/oskca/gopherjs-json"
	"github.com/oskca/gopherjs-vue"
	"honnef.co/go/js/xhr"
	"strconv"
)

//go:generate bash ./makejs.sh

type FrontModel struct {
	*js.Object
	DispPrj           bool                 `js:"DispPrj"`
	TextFilter        string               `js:"textfilter"`
	Projects          []*fm.Project        `js:"projects"`
	SortList          []*fm.SortCol        `js:"sortlist"`
	ColFilterGroup    *fm.ColFilterGroup   `js:"colfilters"`
	EditedPrj         *fm.Project          `js:"editedprj"`
	EditedPrjStat     *fm.ProjectStat      `js:"editedprjstat"`
	PrjStatSignatures *fm.ProjectStatNames `js:"prjstatsignatures"`
	Statuts           []*fm.ValText        `js:"statuts"`
	Types             []*fm.ValText        `js:"types"`
	Risks             []*fm.ValText        `js:"risks"`
	MilestoneKeys     []*fm.ValText        `js:"milestonekeys"`
}

func NewFrontModel(msg string) *FrontModel {
	m := &FrontModel{Object: js.Global.Get("Object").New()}
	m.TextFilter = ""
	m.DispPrj = false
	m.Projects = nil
	m.SortList = []*fm.SortCol{
		fm.NewSortCol("Client", true),
		fm.NewSortCol("Projet", true),
	}
	m.initColFilter()
	m.EditedPrj = fm.NewProject()
	m.EditedPrjStat = nil
	m.PrjStatSignatures = nil
	m.Statuts = biz.CreateStatuts()
	m.Types = biz.CreateTypes()
	m.Risks = biz.CreateRisks()
	m.MilestoneKeys = biz.CreateMilestoneKeys()
	return m
}

func (m *FrontModel) initColFilter() {
	m.ColFilterGroup = fm.NewColFilterGroup()
	m.ColFilterGroup.AddColFilter(
		"Status",
		func(p *fm.Project) string {
			return p.Status
		})
	m.ColFilterGroup.AddColFilter(
		"Type",
		func(p *fm.Project) string {
			return p.Type
		})
	m.ColFilterGroup.AddColFilter(
		"Lead Dev",
		func(p *fm.Project) string {
			return p.LeadDev
		})
	m.ColFilterGroup.AddColFilter(
		"PS",
		func(p *fm.Project) string {
			return p.LeadPS
		})
}

func (m *FrontModel) GetPtf() {
	go m.callGetPtf()
}

func (m *FrontModel) callGetPtf() {
	req := xhr.NewRequest("GET", "/ptf")
	req.Timeout = 2000
	req.ResponseType = xhr.JSON
	m.DispPrj = false
	err := req.Send(nil)
	if err != nil {
		println("Req went wrong : ", err, req.Status)
	}
	m.Projects = []*fm.Project{}
	req.Response.Call("forEach", func(item *js.Object) {
		p := fm.ProjectFromJS(item)
		m.Projects = append(m.Projects, p)
	})
	m.DispPrj = true
	//js.Global.Set("resp", m.Resp)
}

func (m *FrontModel) callGetStatPrjList() {
	req := xhr.NewRequest("GET", "/stat/prjlist/"+strconv.Itoa(m.EditedPrj.Id))
	req.Timeout = 2000
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		println("Req went wrong : ", err, req.Status)
	}
	if req.Status == 200 {
		m.PrjStatSignatures = fm.NewProjectStatNameFromJS(req.Response)
	}
	//TODO Manage Status != 200
}

func (m *FrontModel) callGetProjectStat() {
	req := xhr.NewRequest("GET", "/stat/"+strconv.Itoa(m.EditedPrj.Id))
	req.Timeout = 1000
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		println("Req went wrong : ", err, req.Status)
		return
	}
	if req.Status == 200 {
		m.EditedPrjStat = fm.NewProjectStatFromJS(req.Response)
	}
	//TODO Manage Status != 200
}

func (m *FrontModel) callUpdatePrj(uprj *fm.Project) {
	req := xhr.NewRequest("PUT", "/ptf/"+strconv.Itoa(m.EditedPrj.Id))
	req.Timeout = 1000
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uprj))
	if err != nil {
		println("Req went wrong : ", err, req.Status)
	}
	if req.Status == 200 {
		m.EditedPrj.Copy(fm.ProjectFromJS(req.Response))
	}
}

func (m *FrontModel) callCreatePrj(uprj *fm.Project) {
	req := xhr.NewRequest("POST", "/ptf")
	req.Timeout = 1000
	req.ResponseType = xhr.JSON
	err := req.Send(json.Stringify(uprj))
	if err != nil {
		println("Req went wrong : ", err, req.Status)
	}
	if req.Status == 201 {
		m.EditedPrj.Copy(fm.ProjectFromJS(req.Response))
	}
	m.Projects = append(m.Projects, m.EditedPrj)
}

func (m *FrontModel) callDeletePrj(dprj *fm.Project) {
	req := xhr.NewRequest("DELETE", "/ptf/"+strconv.Itoa(m.EditedPrj.Id))
	req.Timeout = 1000
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		println("Req went wrong : ", err, req.Status)
	}
	if req.Status == 200 {
		m.deletePrj(dprj)
	}
}

func (m *FrontModel) deletePrj(dprj *fm.Project) {
	for i, p := range m.Projects {
		if p.Id == dprj.Id {
			m.Projects = append(m.Projects[:i], m.Projects[i+1:]...)
			break
		}
	}
}

func (m *FrontModel) ProcessEditedPrj() {
	if m.EditedPrj.Id >= 0 {
		go m.callUpdatePrj(m.EditedPrj)
	} else {
		go m.callCreatePrj(m.EditedPrj)
	}
}

func (m *FrontModel) DeleteEditedPrj() {
	if m.EditedPrj.Id >= 0 {
		go m.callDeletePrj(m.EditedPrj)
	}
}

func (m *FrontModel) EditProject(p *fm.Project) {
	m.EditedPrj = p
	m.showEditProjectModal()
}

func (m *FrontModel) ShowProjectStat(p *fm.Project) {
	m.EditedPrj = p
	go func() {
		m.callGetProjectStat()
		m.showProjectStatModal(p)
	}()
}

func (m *FrontModel) RefreshColFilter() {
	println("TODO process RefreshColFilter")
	myVM.Get("colfilteredprojects")
}

func (m *FrontModel) CreateNewProject() {
	m.EditedPrj = fm.NewProject()
	m.showEditProjectModal()
}

func (m *FrontModel) showEditProjectModal() {
	go func() {
		m.callGetStatPrjList()
		jQuery("#EditProjectModalComp").Get(0).Get("__vue__").Call("ShowEditProjectModal", m.EditedPrj, m.PrjStatSignatures)
	}()
}

func (m *FrontModel) showProjectStatModal(p *fm.Project) {
	jQuery("#ProjectStatModalComp").Get(0).Get("__vue__").Call("ShowProjectStatModal", p, m.EditedPrjStat)
}

func (m *FrontModel) IsDisplayed(p *fm.Project) bool {
	return fm.TextFiltered(p, m.TextFilter) && m.ColFilterGroup.ColFiltered(p)
}

func (m *FrontModel) RiskIcon(p *fm.Project) string {
	var res string
	switch p.Risk {
	case "2":
		res = "red warning sign icon"
	case "1":
		res = "orange warning circle icon"
	default:
		//res = "green info circle icon"
	}
	return res
}

var jQuery = jquery.NewJQuery
var myVM *vue.ViewModel = nil

func main() {
	comps.RegisterEditProjectModalComp()
	comps.RegisterColTitleComp()
	comps.RegisterColTitleWithFilterComp()
	comps.RegisterDateTableCellComp()
	comps.RegisterWorkLoadCellComp()
	comps.RegisterProjectStatModalComp()

	comps.ChargeFilterRegister("ChargeFormat")
	comps.DateFilterRegister("DateFormat")

	mo := vue.NewOption()
	mo.El = "#prj-app"
	mo.SetDataWithMethods(NewFrontModel(""))

	mo.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		m := &FrontModel{Object: vm.Object}
		m.GetPtf()
	})

	mo.AddComputed("sortedProjects", func(vm *vue.ViewModel) interface{} {
		m := &FrontModel{Object: vm.Object}
		sortCols := m.SortList
		m.Projects = fm.SortedProjects(m.Projects, sortCols)
		//println("Did Sorting ...")
		return m.Projects
	})

	//vm := vue.New("#prj-app", NewFrontModel(""))
	myVM = mo.NewViewModel()

	js.Global.Set("vm", myVM)
}

// Done gauge on WorkLoad Column (create new component)
// TODO Manage Different Columns Sets
