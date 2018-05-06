package projecttree

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/novagile/src/client/tools"
)

var id int

type Node struct {
	*js.Object

	Id       int     `js:"id"`
	Label    string  `js:"label"`
	Children []*Node `js:"children"`
	Parent   *Node   `js:"parent"`

	HRef    string `js:"href"`
	Summary string `js:"summary"`

	ParentRatio  bool    `js:"parentRatio"`
	LinkedHour   float64 `js:"lhour"`
	UnlinkedHour float64 `js:"nlhour"`
}

func NewNode(label string) *Node {
	n := &Node{Object: tools.O()}
	n.Id = id
	id++
	n.Label = label
	n.Children = []*Node{}
	n.Parent = nil
	n.HRef = ""
	n.Summary = ""
	n.ParentRatio = false
	n.LinkedHour = 0
	n.UnlinkedHour = 0
	return n
}

func (n *Node) AddChild(c *Node) {
	n.Children = append(n.Children, c)
	c.Parent = n
}

func (n *Node) SetIssueInfo(issue, summary string) {
	n.HRef = tools.UrlJiraBrowseIssue + issue
	n.Summary = summary
}

func (n *Node) SetHour(lot string, hour float64) {
	if lot == "" {
		n.UnlinkedHour = hour
		return
	}
	n.LinkedHour = hour
}

func (n *Node) Update() {
	if len(n.Children) == 0 {
		return
	}

	for _, cn := range n.Children {
		cn.Update()
		n.LinkedHour += cn.LinkedHour
		n.UnlinkedHour += cn.UnlinkedHour
	}

	compHour := func(ci, cj *Node) int {
		a := ci.LinkedHour + ci.UnlinkedHour
		b := cj.LinkedHour + cj.UnlinkedHour
		if a < b {
			return 1
		}
		if a > b {
			return -1
		}
		return 0
	}

	compRatio := func(ci, cj *Node) int {
		a := 0.0
		if (ci.LinkedHour + ci.UnlinkedHour) > 0 {
			a = ci.LinkedHour / (ci.LinkedHour + ci.UnlinkedHour)
		}
		b := 0.0
		if (cj.LinkedHour + cj.UnlinkedHour) > 0 {
			b = cj.LinkedHour / (cj.LinkedHour + cj.UnlinkedHour)
		}
		if a > b {
			return 1
		}
		if a < b {
			return -1
		}
		return 0
	}

	if n.Children[0].ParentRatio {
		n.Object.Get("children").Call("sort", compHour)
	} else {
		n.Object.Get("children").Call("sort", compRatio)
	}
}
