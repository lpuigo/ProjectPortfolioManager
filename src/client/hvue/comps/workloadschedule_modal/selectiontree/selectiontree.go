package selectiontree

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wsr "github.com/lpuig/novagile/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/novagile/src/client/tools"
)

const template = `
<el-tree
	ref="tree"
    :data="nodes"
    :props="nodeProps"
	node-key="id"
	:default-checked-keys="checkedNodes"
	:render-after-expand="false"
    show-checkbox
	accordion
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
	nl := []int{}
	stcm.Nodes = []*Node{}

	nl = append(nl, stcm.createNodeProject()...)
	nl = append(nl, stcm.createNodeByCrit(
		NewNode("By Lead Dev", nil),
		func(r *wsr.WorkloadScheduleRecord) string {
			return r.LeadDev
		})...)
	nl = append(nl, stcm.createNodeByCrit(
		NewNode("By Lead PS", nil),
		func(r *wsr.WorkloadScheduleRecord) string {
			return r.LeadPS
		})...)
	nl = append(nl, stcm.createNodeByCrit(
		NewNode("By Status", nil),
		func(r *wsr.WorkloadScheduleRecord) string {
			return r.Status
		})...)
	//nl = append(nl, stcm.createNodeProject()...)

	stcm.CheckedNodes = nl
}

func (stcm *SelectionTreeCompModel) createNodeProject() (list []int) {
	projectNode := NewNode("Projects List", nil)
	for _, wr := range stcm.WrkSched.Records {
		wr.Display = true
		n := NewNode(wr.Name, wr)
		projectNode.AddChild(n)
		list = append(list, n.Id)
	}
	stcm.Nodes = append(stcm.Nodes, projectNode)
	return
}

type groupBy func(r *wsr.WorkloadScheduleRecord) string

func (stcm *SelectionTreeCompModel) createNodeByCrit(parent *Node, crit groupBy) (list []int) {
	nodesBy := map[string]*Node{}

	for _, wr := range stcm.WrkSched.Records {
		chldNode := NewNode(wr.Name, wr)
		curCrit := crit(wr)
		if curCrit == "" {
			curCrit = "<Blank>"
		}
		critNode, found := nodesBy[curCrit]
		if !found {
			critNode = NewNode(curCrit, nil)
			nodesBy[curCrit] = critNode
			parent.AddChild(critNode)
		}
		critNode.AddChild(chldNode)
		list = append(list, chldNode.Id)
	}
	stcm.Nodes = append(stcm.Nodes, parent)
	return
}

func (stcm *SelectionTreeCompModel) HandleCheckChange(node *Node, checked, indeterminate bool) {
	if node.WrkSchedRec.Object == nil {
		return
	}
	node.WrkSchedRec.Display = checked

	stcm.updateNodesCheckState()
	//stcm.VM.Emit("update:wrkSched", stcm.WrkSched)
}

func (stcm *SelectionTreeCompModel) updateNodesCheckState() {
	println("updateNodesCheckState")
	selectedNodes := []int{}
	for _, n := range stcm.Nodes {
		selectedNodes = append(selectedNodes, n.updateCheckState()...)
	}
	stcm.VM.Refs("tree").Call("setCheckedKeys", selectedNodes)
}
