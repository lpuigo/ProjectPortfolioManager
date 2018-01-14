package Comps

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/oskca/gopherjs-vue"
	"strconv"
)

type IssueChartComp struct {
	*js.Object
	IssueStat *fm.IssueStat `js:"issuestat"`
	Height    float64       `js:"height"`
	Width     float64       `js:"width"`
	XTicks    []float64     `js:"xticks"`
}

func NewIssueChartComp() *IssueChartComp {
	ic := &IssueChartComp{Object: js.Global.Get("Object").New()}
	ic.IssueStat = nil
	ic.Height = 70
	ic.Width = 0
	ic.XTicks = []float64{}
	return ic
}

const (
	xMinOffset  = 15.0
	xMaxOffset  = 15.0
	yMinOffset  = 15.0
	yMaxOffset  = 15.0
	ticksLenght = 5.0
)

// RegisterIssueChartComp registers to current vue intance a IssueChart component
// having the following profile
//  <issue-chart :issuestat="some issuestat"></issue-chart>
func RegisterIssueChartComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewIssueChartComp

	o.AddProp("issuestat")

	o.Template = `
	<div>
		<h5 v-if="issuestat.href">Evolution des temps : <a :href="issuestat.href" target="_blank">{{issuestat.issue}}</a></h5>
		<h5 v-else>Evolution des temps : {{issuestat.issue}}</h5>
		<svg :height="height" style="background: whitesmoke; width: 100%">
			<polyline :points="XAxisPoints()" style="fill:none;stroke:#d7d7d7;stroke-width:1" />
			<polyline :points="YAxisPoints()" style="fill:none;stroke:#d7d7d7;stroke-width:1" />
			<!--<g v-for="(date, i) in issuestat.dates">-->
				<!--<polyline :points="XTicksPoints(i)" style="fill:none;stroke:#d7d7d7;stroke-width:1" />-->
			<!--</g>-->
			<polyline :points="StrokeLinePoints(issuestat.timeEstimated)" style="fill:none;stroke:red;stroke-width:2" />
			<polyline :points="StrokeLinePoints(issuestat.timeRemaining)" style="fill:none;stroke:orange;stroke-width:2" />
			<polyline :points="StrokeLinePoints(issuestat.timeSpent)" style="fill:none;stroke:teal;stroke-width:4" />
		</svg>
	</div>`

	o.AddMethodWithReturnValue("XAxisPoints", func(vm *vue.ViewModel, args []*js.Object) interface{} {
		ic := &IssueChartComp{Object: vm.Object}
		y := ic.Height - yMinOffset
		x1 := xMinOffset - ticksLenght
		x2 := ic.Width - xMaxOffset

		return pointCoord(x1, y) + pointCoord(x2, y)
	})

	o.AddMethodWithReturnValue("YAxisPoints", func(vm *vue.ViewModel, args []*js.Object) interface{} {
		ic := &IssueChartComp{Object: vm.Object}
		x := xMinOffset
		y1 := ic.Height - yMinOffset + ticksLenght
		y2 := yMaxOffset

		return pointCoord(x, y1) + pointCoord(x, y2)
	})

	o.AddMethodWithReturnValue("XTicksPoints", func(vm *vue.ViewModel, args []*js.Object) interface{} {
		ic := &IssueChartComp{Object: vm.Object}
		println("XTicksPoints", ic.Object)
		i := args[0].Int()

		x := ic.XTicks[i]
		y1 := ic.Height - yMinOffset
		y2 := y1 + ticksLenght

		return pointCoord(x, y1) + pointCoord(x, y2)
	})

	// returns strokeline points string for given serie
	o.AddMethodWithReturnValue("StrokeLinePoints", func(vm *vue.ViewModel, args []*js.Object) interface{} {
		ic := &IssueChartComp{Object: vm.Object}

		ys := args[0].Interface().([]float64)
		nbPoints := len(ys)
		trX, trY := ic.genTranslateFunc(float64(nbPoints), 50)
		res := ""
		for i := 0; i < nbPoints; i++ {
			x := trX(float64(i) + 0.5)
			y := trY(ys[i])
			res += pointCoord(x, y)
		}
		return res
	})

	o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		ic := &IssueChartComp{Object: vm.Object}
		ic.Width = vm.El.Call("getBoundingClientRect").Get("width").Float()
		ic.genXTicks()
		if ic.IssueStat.HRef == "" {
			ic.Height = 130 // Project Overview, make it bigger
		}
		println("Mounted", ic.Object)
	})

	return o.NewComponent().Register("issue-chart")
}

func (ic *IssueChartComp) genXTicks() {
	nbTicks := len(ic.IssueStat.Dates)
	trX, _ := ic.genTranslateFunc(float64(nbTicks), 50)
	tx := []float64{}
	for i := 0; i < nbTicks; i++ {
		tx = append(tx, trX(float64(i)+0.5))
	}
	ic.XTicks = tx
	println("genXTicks", ic.Object)
}

func (ic *IssueChartComp) genTranslateFunc(xMax, yMax float64) (trX, trY func(float64 float64) float64) {
	xratio := (ic.Width - xMinOffset - xMaxOffset) / xMax
	yratio := (ic.Height - yMinOffset - yMaxOffset) / yMax
	trX = func(x float64) float64 {
		return x*xratio + xMinOffset
	}
	trY = func(y float64) float64 {
		return ic.Height - yMinOffset - y*yratio
	}
	return
}

func pointCoord(x, y float64) string {
	return strconv.FormatFloat(x, 'f', 1, 64) + "," + strconv.FormatFloat(y, 'f', 1, 64) + " "
}
