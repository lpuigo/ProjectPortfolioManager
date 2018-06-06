package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/prjptf/src/client/frontmodel"
	"github.com/lpuig/prjptf/src/client/tools"
)

type TimeLine struct {
	*js.Object
	Name       string            `js:"name"`
	Project    *fm.Project       `js:"project"`
	Phases     []*Phase          `js:"phases"`
	MileStones map[string]string `js:"milestones"`
}

func NewTimeLine(p *fm.Project) *TimeLine {
	t := &TimeLine{Object: tools.O()}
	t.Project = p
	t.Name = p.Client + " - " + p.Name
	t.Phases = []*Phase{}
	t.MileStones = p.MileStones
	return t
}

func (t *TimeLine) AddPhase(p *Phase) {
	t.Phases = append(t.Phases, p)
}
