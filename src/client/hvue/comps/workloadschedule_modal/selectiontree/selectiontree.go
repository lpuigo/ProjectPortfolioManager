package selectiontree

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wsr "github.com/lpuig/novagile/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/novagile/src/client/tools"
)

const template = `
<el-tree
    :data="nodes"
    :props="nodeProps"
	node-key="id"
	:default-checked-keys="checkedNodes"
	:render-after-expand="false"
    show-checkbox
	@check-change="HandleCheckChange"
>
</el-tree>
`

func Register() {
	hvue.NewComponent("selection-tree",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		//hvue.Props("selection"),
		hvue.Template(template),
		hvue.Props("wrkSched"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewSelectionTreeCompModel(vm)
		}),
		hvue.MethodsOf(&SelectionTreeCompModel{}),
	}
}

type SelectionTreeCompModel struct {
	*js.Object

	WrkSched     *wsr.WorkloadSchedule `js:"wrkSched"`
	Nodes        []*Node               `js:"nodes"`
	NodeProps    js.M                  `js:"nodeProps"`
	CheckedNodes []int                 `js:"checkedNodes"`

	VM *hvue.VM `js:"VM"`
}

func NewSelectionTreeCompModel(vm *hvue.VM) *SelectionTreeCompModel {
	stcm := &SelectionTreeCompModel{Object: tools.O()}
	stcm.NodeProps = GetNodeProps()
	stcm.Nodes = []*Node{}
	stcm.CheckedNodes = []int{}

	stcm.VM = vm
	return stcm
}

func (stcm *SelectionTreeCompModel) Init(wsr *wsr.WorkloadSchedule) {
	stcm.WrkSched = wsr
	stcm.updateNodes()
}

func (stcm *SelectionTreeCompModel) updateNodes() {
	projectNode := NewNode("Projects", nil)
	stcm.Nodes = []*Node{projectNode}
	nl := []int{}
	for _, wr := range stcm.WrkSched.Records {
		wr.Display = true
		n := NewNode(wr.Name, wr)
		projectNode.AddChild(n)
		nl = append(nl, n.Id)
	}
	stcm.CheckedNodes = nl
}

func (stcm *SelectionTreeCompModel) HandleCheckChange(node *Node, checked, indeterminate bool) {
	if node.WrkSchedRec.Object == nil {
		return
	}
	node.WrkSchedRec.Display = checked
	stcm.VM.Emit("update:wrkSched", stcm.WrkSched)
}
