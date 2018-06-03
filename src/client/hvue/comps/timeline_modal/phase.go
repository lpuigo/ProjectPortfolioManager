package timeline_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/prjptf/src/client/tools"
	"strconv"
)

type Phase struct {
	*js.Object
	Comment string  `js:"comment"`
	Class   string  `js:"class"`
	Style   string  `js:"style"`
	IsFirst bool    `js:"first"`
	IsLast  bool    `js:"last"`
	Offset  float64 `js:"offset"`
	Length  float64 `js:"lenght"`
}

func NewPhase(comment, class string) *Phase {
	p := &Phase{Object: tools.O()}
	p.Comment = comment
	p.Class = class
	p.Style = ""
	p.IsFirst = false
	p.IsLast = false
	return p
}

func (p *Phase) SetPositionInfo(offset, length float64) {
	p.Offset = offset
	p.Length = length
}

func (p *Phase) SetStyleClass(currentOffset float64) float64 {
	s := ""
	switch {
	case p.Offset > 0:
		s += "margin-left: " + strconv.FormatFloat(p.Offset, 'f', 1, 64) + "%; "
	case p.Offset < 0:
		p.Length += p.Offset
		p.Offset = 0
		p.IsFirst = false
	}
	if currentOffset+p.Offset+p.Length > 100 {
		p.Length = 100 - p.Offset - currentOffset
		p.IsLast = false
	}
	s += "width: " + strconv.FormatFloat(p.Length, 'f', 1, 64) + "%"
	p.Style = s

	if p.IsFirst {
		p.Class += " first"
	}
	if p.IsLast {
		p.Class += " last"
	}
	return currentOffset+p.Offset+p.Length
}
