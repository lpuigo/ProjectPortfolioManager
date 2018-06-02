package selectiontree

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wsr "github.com/lpuig/prjptf/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/prjptf/src/client/tools"
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
	@check="HandleCheck"
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
	projectNode.sortChildren(byLabel)
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
	parent.sortChildren(byLabel)
	stcm.Nodes = append(stcm.Nodes, parent)
	return
}

func (stcm *SelectionTreeCompModel) HandleCheck(node *Node, obj *js.Object) {
	checked := obj.Get("checkedKeys").Call("includes", node.Id).Bool()

	stcm.handleCheck(node, checked)
	stcm.VM.Emit("update:wrkSched", stcm.WrkSched)
}

func (stcm *SelectionTreeCompModel) handleCheck(node *Node, checkstate bool) {
	if node.WrkSchedRec.Object == nil { // Non Leaf Node
		for _, cn := range node.Children {
			stcm.handleCheck(cn, checkstate)
		}
		return
	}

	// Leaf Node
	if node.WrkSchedRec.Display == checkstate { // already in correct state, skip
		return
	}

	node.WrkSchedRec.Display = checkstate
	this := stcm.VM.Refs("tree")
	for _, nodeId := range stcm.getRelatedNodeId(node.WrkSchedRec.Id) {
		if node.Id == nodeId {
			continue
		}
		this.Call("setChecked", nodeId, node.WrkSchedRec.Display, false)
	}
}

func (stcm *SelectionTreeCompModel) getRelatedNodeId(projectId int) (list []int) {
	for _, n := range stcm.Nodes {
		list = append(list, n.getRelatedNodeId(projectId)...)
	}
	return
}
