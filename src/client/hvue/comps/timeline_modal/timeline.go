package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/novagile/src/client/tools"
)

type TimeLine struct {
	*js.Object
	Name   string   `js:"name"`
	Phases []*Phase `js:"phase"`
}

func NewTimeLine(name string) *TimeLine {
	t := &TimeLine{Object: tools.O()}
	t.Name = name
	t.Phases = []*Phase{}
	return t
}

func (t *TimeLine) AddPhase(p *Phase) {
	t.Phases = append(t.Phases, p)
}
