package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/novagile/src/client/tools"
	"strconv"
)

type Phase struct {
	*js.Object
	Name  string `js:"name"`
	Style string `js:"style"`
}

func NewPhase(name string) *Phase {
	p := &Phase{Object: tools.O()}
	p.Name = name
	p.Style = ""
	return p
}

func (p *Phase) SetStyle(offset, length float64) {
	s := ""
	if offset > 0 {
		s += "margin-left: " + strconv.FormatFloat(offset, 'f', 1, 64) + "%; "
	}
	s += "width: " + strconv.FormatFloat(length, 'f', 1, 64) + "%"
	p.Style = s
}
