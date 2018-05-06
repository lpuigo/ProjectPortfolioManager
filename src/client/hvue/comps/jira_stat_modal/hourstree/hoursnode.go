package hourstree

import "github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/node"

type Node struct {
	*node.Node
	Name    string    `js:"name"`
	Hours   []float64 `js:"hours"`
	MaxHour float64   `js:"maxHour"`
}

func NewNode(label string, hours []float64, maxhour float64) *Node {
	n := &Node{Node: node.New(label, nil)}
	n.Hours = hours
	n.MaxHour = maxhour
	return n
}

func (hn *Node) AddChild(c *Node) {
	hn.Node.AddChild(c)

	if len(hn.Hours) == 0 {
		hn.Hours = make([]float64, len(c.Hours))
		hn.MaxHour = 0
	}

	hn.MaxHour += c.MaxHour

	for i, h := range c.Hours {
		//hn.Hours[i] += h // Gopherjs style :(
		o := hn.Object.Get("hours")
		o.SetIndex(i, o.Index(i).Float()+h)
	}
}
