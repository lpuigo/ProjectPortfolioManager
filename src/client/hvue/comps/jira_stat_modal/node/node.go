package node

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/prjptf/src/client/tools"
)

var id int

type NodeChild interface {
}

type Node struct {
	*js.Object

	Id       int     `js:"id"`
	Label    string  `js:"label"`
	Children []NodeChild `js:"children"`
}

func New(label string, children []NodeChild) *Node {
	n := &Node{Object: tools.O()}
	n.Id = id
	n.Label = label
	n.Children = children
	id++
	return n
}

func (n *Node) AddChild(c NodeChild) {
	n.Children = append(n.Children, c)
}
