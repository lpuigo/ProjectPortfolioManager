package node

type HoursNode struct {
	*Node
	Name    string    `js:"name"`
	Hours   []float64 `js:"hours"`
	MaxHour float64   `js:"maxHour"`
}

func NewHoursNode(label string, hours []float64, maxhour float64) *HoursNode {
	n := &HoursNode{Node: NewNode(label, nil)}
	n.Hours = hours
	n.MaxHour = maxhour
	return n
}

func (hn *HoursNode) AddChild(c *HoursNode) {
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
