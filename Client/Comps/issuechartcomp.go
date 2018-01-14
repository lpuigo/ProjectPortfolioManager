package Comps

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/oskca/gopherjs-vue"
	"strconv"
)

type IssueChartComp struct {
	*js.Object
	PrjStat *fm.ProjectStat `js:"projectstat"`
	Height  float64         `js:"height"`
	Width   float64         `js:"width"`
}

func NewIssueChartComp() *IssueChartComp {
	ic := &IssueChartComp{Object: js.Global.Get("Object").New()}
	ic.PrjStat = nil
	ic.Height = 100
	return ic
}

// RegisterIssueChartComp registers to current vue intance a IssueChart component
// having the following profile
//  <issue-chart :projectstat="some projectstat"></issue-chart>
func RegisterIssueChartComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewIssueChartComp

	o.AddProp("projectstat")

	o.Template = `
	<div style="background: gray">
		<svg :height="height" style="background: whitesmoke; width: 100%" v-if="projectstat">
			<g v-for="(issue, index) in projectstat.issues">
				<polyline :points="StrokeLinePoints(issue, index)" style="fill:none;stroke:black;stroke-width:3" />
			</g>
		</svg>
		<pre>Stats : {{projectstat}}</pre>
	</div>`

	o.AddMethod("RenderChart", func(vm *vue.ViewModel, args []*js.Object) {
		ic := &IssueChartComp{Object: vm.Object}
		ic.PrjStat = &fm.ProjectStat{Object: args[0]}
		ic.Height = 50
		ic.Width = vm.El.Call("getBoundingClientRect").Get("width").Float()
	})

	// returns strokeline points string for given serie
	o.AddMethodWithReturnValue("StrokeLinePoints", func(vm *vue.ViewModel, args []*js.Object) interface{} {
		ic := &IssueChartComp{Object: vm.Object}
		ys := ic.PrjStat.TimeSpent[args[1].Int()]
		nbPoints := len(ic.PrjStat.Dates)
		trX, trY := ic.genTranslateFunc(float64(nbPoints+2), 50)
		res := ""
		for i := 0; i < nbPoints; i++ {
			x := trX(float64(i) + 1.0)
			y := trY(ys[i] + 2)
			res += strconv.FormatFloat(x, 'f', 1, 64) + "," + strconv.FormatFloat(y, 'f', 1, 64) + " "
		}
		return res
	})

	return o.NewComponent().Register("issue-chart")
}

func (ic *IssueChartComp) genTranslateFunc(xMax, yMax float64) (trX, trY func(float64 float64) float64) {
	xratio := ic.Width / xMax
	yratio := ic.Height / yMax
	trX = func(x float64) float64 {
		return x * xratio
	}
	trY = func(y float64) float64 {
		return ic.Height - y*yratio
	}
	return
}
