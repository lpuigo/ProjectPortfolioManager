package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/novagile/src/client/tools"
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
