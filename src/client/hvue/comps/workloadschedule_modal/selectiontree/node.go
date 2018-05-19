package selectiontree

import (
	"github.com/gopherjs/gopherjs/js"
	wsr "github.com/lpuig/novagile/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/novagile/src/client/tools"
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

func GetNodeProps() js.M {
	return js.M{
		"children": "children",
		"label":    "label",
	}
}
