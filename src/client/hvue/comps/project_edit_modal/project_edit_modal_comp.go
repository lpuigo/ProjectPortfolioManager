package project_edit_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/element/model"
	"github.com/lpuig/novagile/src/client/business"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"time"
)

func Register() {
	hvue.NewComponent("project-edit-modal",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("edited_project"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewProjectEditModalModel(vm)
		}),
		hvue.MethodsOf(&ProjectEditModalModel{}),
		hvue.Computed("unusedMilestoneKeys", func(vm *hvue.VM) interface{} {
			m := &ProjectEditModalModel{Object: vm.Object}
			keyList := []string{}
			vm.Get("milestonesList").Call("forEach", func(vt *js.Object) {
				k := vt.Get("value").String()
				if _, ok := m.CurrentProject.MileStones[k]; ok == false {
					keyList = append(keyList, k)
				}
			})
			return keyList
		}),
		hvue.Computed("usedMilestoneKeys", func(vm *hvue.VM) interface{} {
			m := &ProjectEditModalModel{Object: vm.Object}
			keyList := []string{}
			vm.Get("milestonesList").Call("forEach", func(vt *js.Object) {
				k := vt.Get("value").String()
				if _, ok := m.CurrentProject.MileStones[k]; ok == true {
					keyList = append(keyList, k)
				}
			})
			return keyList
		}),
	}
}

type ProjectEditModalModel struct {
	*js.Object

	EditedProject  *fm.Project `js:"edited_project"`
	CurrentProject *fm.Project `js:"currentProject"`

	Visible           bool `js:"visible"`
	IsNewProject      bool `js:"isNewProject"`
	ShowConfirmDelete bool `js:"showconfirmdelete"`

	RiskList       []*fm.ValText `js:"riskList"`
	StatusList     []*fm.ValText `js:"statusList"`
	TypeList       []*fm.ValText `js:"typeList"`
	MilestonesList []*fm.ValText `js:"milestonesList"`

	ClientNameLookUp bool          `js:"clientNameLookup"`
	ClientNameList   []*fm.ValText `js:"clientNameList"`

	VM *hvue.VM `js:"VM"`
}

func NewProjectEditModalModel(vm *hvue.VM) *ProjectEditModalModel {
	pemm := &ProjectEditModalModel{Object: model.O()}
	pemm.EditedProject = fm.NewProject()
	pemm.CurrentProject = fm.NewProject()
	pemm.Visible = false
	pemm.IsNewProject = false
	pemm.ShowConfirmDelete = false
	pemm.RiskList = business.CreateRisks()
	pemm.StatusList = business.CreateStatuts()
	pemm.TypeList = business.CreateTypes()
	pemm.MilestonesList = business.CreateMilestoneKeys()
	pemm.ClientNameLookUp = false
	pemm.ClientNameList = nil
	pemm.VM = vm
	return pemm
}

func (pemm *ProjectEditModalModel) Show(p *fm.Project) {
	pemm.EditedProject = p
	pemm.CurrentProject = fm.NewProject()
	pemm.CurrentProject.Copy(pemm.EditedProject)
	pemm.IsNewProject = false
	pemm.ShowConfirmDelete = false
	pemm.ClientNameList = nil
	pemm.ClientNameLookUp = false
	pemm.Visible = true
}

func (pemm *ProjectEditModalModel) Hide() {
	pemm.Visible = false
	pemm.ShowConfirmDelete = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Milestones Button Methods

func (pemm *ProjectEditModalModel) DeleteMilestone(ms string) {
	//pemm = &ProjectEditModalModel{Object: vm.Object}
	pemm.CurrentProject.RemoveMileStone(ms)
}

func (pemm *ProjectEditModalModel) AddMilestone(ms string) {
	//pemm = &ProjectEditModalModel{Object: vm.Object}
	pemm.CurrentProject.AddMileStone(ms)
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (pemm *ProjectEditModalModel) ConfirmChange() {
	pemm.EditedProject.Copy(pemm.CurrentProject)
	pemm.VM.Emit("update:edited_project", pemm.EditedProject)
	pemm.Hide()
}

func (pemm *ProjectEditModalModel) DeleteProject() {
	pemm.VM.Emit("delete:edited_project", pemm.EditedProject)
	pemm.Hide()
}

func (pemm *ProjectEditModalModel) Duplicate() {
	pemm.EditedProject = pemm.CurrentProject
	pemm.CurrentProject.Name += " (Copy)"
	pemm.CurrentProject.Id = -1
	pemm.CurrentProject.CurrentWL = 0.0
	pemm.IsNewProject = true
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Get Client-Name List

func (pemm *ProjectEditModalModel) SetClientName(vt *fm.ValText) {
	pemm.CurrentProject.Client = vt.Value
	pemm.CurrentProject.Name = vt.Text
}

func (pemm *ProjectEditModalModel) HasClientNameList() bool {
	return len(pemm.ClientNameList) > 0
}

func (pemm *ProjectEditModalModel) GetClientNameList() {
	if pemm.HasClientNameList() {
		return
	}
	pemm.ClientNameLookUp = true
	go pemm.callClientNameList()
}

func (pemm *ProjectEditModalModel) callClientNameList() {
	// TODO Implement XHR query to retrieve PrjStatList
	time.Sleep(2 * time.Second)
	pemm.ClientNameList = []*fm.ValText{
		fm.NewValText("Client A", "Prj A1"),
		fm.NewValText("Client A", "Prj A2"),
		fm.NewValText("Client B", "Prj B1"),
		fm.NewValText("Client B", "Prj B2"),
	}
	pemm.ClientNameLookUp = false
}
