package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/novagile/src/client/tools"
)

type TimeLine struct {
	*js.Object
	Name  string `js:"name"`
	Style string `js:"style"`
}

func NewTimeLine(name string) *TimeLine {
	p := &TimeLine{Object: tools.O()}
	p.Name = name
	p.Style = ""
	return p
}

func (p *TimeLine) SetStyle() {
	p.Style = ""
}
