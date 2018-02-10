package comps

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"strconv"
	"time"
)

type DateCellComp struct {
	*js.Object
	Date string `js:"date"`
}

func NewDateCellComp() *DateCellComp {
	dc := &DateCellComp{Object: js.Global.Get("Object").New()}
	dc.Date = ""
	return dc
}

func (dc *DateCellComp) getDate() time.Time {
	d, err := time.Parse("2006-01-02", dc.Date)
	if err != nil {
		return time.Time{}
	}
	return d
}

// CompareWeek returns -1 if given dc is at least one week before today, 0 if in same week, 1 if at least one week later
func (dc *DateCellComp) CompareWeek() int {
	formatweek := func(y, w int) string {
		if w < 10 {
			return strconv.Itoa(y) + "0" + strconv.Itoa(w)
		}
		return strconv.Itoa(y) + strconv.Itoa(w)
	}

	dd := dc.getDate()
	if dd.IsZero() {
		return 1
	}

	nw := formatweek(time.Now().ISOWeek())
	dw := formatweek(dd.ISOWeek())
	switch {
	case dw < nw:
		return -1
	case dw > nw:
		return 1
	default:
		return 0
	}
}

// RegisterDateTableCellComp registers to current vue intance a DateTableCell component
// having the following profile
//  <td is="date-cell" :date="some date"></td>
func RegisterDateTableCellComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewDateCellComp

	o.AddProp("date")

	//<td class="collapsing center aligned disabled">{{prj.milestones.Kickoff | DateFormat}}</td>
	o.Template = `
	<td class="collapsing center aligned" :class="displaymode">
	{{date | DateFormat}}
	</td>`

	o.AddComputed("displaymode", func(vm *vue.ViewModel) interface{} {
		ct := &DateCellComp{Object: vm.Object}
		switch ct.CompareWeek() {
		case -1:
			return "positive"
		case 0:
			return "warning"
		default:
			return ""
		}
	})

	return o.NewComponent().Register("date-cell")
}
