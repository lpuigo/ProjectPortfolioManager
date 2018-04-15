package project_edit_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/element/model"
	"github.com/lpuig/novagile/src/client/business"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
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
	}
}

type ProjectEditModalModel struct {
	*js.Object

	EditedProject  *fm.Project `js:"edited_project"`
	CurrentProject *fm.Project `js:"currentProject"`

	Visible           bool `js:"visible"`
	IsNewProject      bool `js:"isNewProject"`
	ShowConfirmDelete bool `js:"showconfirmdelete"`

	RiskList   []*fm.ValText `js:"riskList"`
	StatusList []*fm.ValText `js:"statusList"`
	TypeList   []*fm.ValText `js:"typeList"`

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
	pemm.VM = vm
	return pemm
}

func (pemm *ProjectEditModalModel) Show(p *fm.Project) {
	pemm.EditedProject = p
	pemm.CurrentProject = fm.NewProject()
	pemm.CurrentProject.Copy(pemm.EditedProject)
	pemm.IsNewProject = false
	pemm.ShowConfirmDelete = false
	pemm.Visible = true
}

func (pemm *ProjectEditModalModel) ConfirmChange() {
	pemm.EditedProject.Copy(pemm.CurrentProject)
	pemm.VM.Emit("update:edited_project", pemm.EditedProject)
	pemm.Visible = false
}

func (pemm *ProjectEditModalModel) Duplicate() {
	pemm.EditedProject = pemm.CurrentProject
	pemm.CurrentProject.Name += " (Copy)"
	pemm.CurrentProject.Id = -1
	pemm.CurrentProject.CurrentWL = 0.0
	pemm.IsNewProject = true
}

func (pemm *ProjectEditModalModel) NewProject() {
	pemm.VM.Emit("update:edited_project", pemm.EditedProject)
	pemm.Visible = false
}

func (pemm *ProjectEditModalModel) DeleteProject() {
	pemm.VM.Emit("delete:edited_project", pemm.EditedProject)
	pemm.ShowConfirmDelete = false
	pemm.Visible = false
}
