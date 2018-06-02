package selectiontree

import (
	"github.com/gopherjs/gopherjs/js"
	wsr "github.com/lpuig/prjptf/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/prjptf/src/client/tools"
)

var id int

type Node struct {
	*js.Object

	Id          int                         `js:"id"`
	Label       string                      `js:"label"`
	WrkSchedRec *wsr.WorkloadScheduleRecord `js:"wrkschedrecord"`
	Children    []*Node                     `js:"children"`
}

func NewNode(label string, rec *wsr.WorkloadScheduleRecord) *Node {
	n := &Node{Object: tools.O()}
	n.Id = id
	id++
	n.Label = label
	n.WrkSchedRec = rec
	n.Children = []*Node{}
	return n
}

func (n *Node) AddChild(c *Node) {
	n.Children = append(n.Children, c)
}

func (n *Node) checkedNodesList() (list []int) {
	if n.WrkSchedRec != nil && n.WrkSchedRec.Object != nil && n.WrkSchedRec.Display {
		list = append(list, n.Id)
	}
	for _, c := range n.Children {
		list = append(list, c.checkedNodesList()...)
	}
	return
}

func (n *Node) updateCheckState() (selected []int) {
	if n.WrkSchedRec == nil || n.WrkSchedRec.Object == nil {
		return
	}
	if n.WrkSchedRec.Display {
		selected = append(selected, n.Id)
	}
	if len(n.Children) == 0 {
		return
	}
	for _, cn := range n.Children {
		selected = append(selected, cn.updateCheckState()...)
	}
	return
}

func (n *Node) getRelatedNodeId(projectId int) (list []int) {
	if n.WrkSchedRec.Object != nil { // its a leaf node
		if n.WrkSchedRec.Id == projectId {
			list = append(list, n.Id)
		}
		return
	}
	for _, cn := range n.Children {
		list = append(list, cn.getRelatedNodeId(projectId)...)
	}
	return
}

func byLabel(a, b *Node) int {
	if a.Label < b.Label {
		return -1
	}
	if a.Label > b.Label {
		return 1
	}
	return 0
}

func (n *Node) sortChildren(comp func(a, b *Node) int) {
	if len(n.Children) == 0 {
		return
	}
	n.Object.Get("children").Call("sort", comp)
	for _, cn := range n.Children {
		cn.sortChildren(comp)
	}
}

func GetNodeProps() js.M {
	return js.M{
		"children": "children",
		"label":    "label",
	}
}
