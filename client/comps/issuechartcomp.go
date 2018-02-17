package comps

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/novagile/client/frontmodel"
	"github.com/oskca/gopherjs-vue"
	"time"
)

const TimeJSLayout string = "2006-01-02"

func dateFromJSString(s string) int64 {
	res := time.Time{}
	res, _ = time.Parse(TimeJSLayout, s)
	return res.Unix() * 1000
}

type IssueChartComp struct {
	*js.Object
	IssueStat *fm.IssueStat `js:"issuestat"`
}

func NewIssueChartComp() *IssueChartComp {
	ic := &IssueChartComp{Object: js.Global.Get("Object").New()}
	ic.IssueStat = nil
	return ic
}

// RegisterIssueChartComp registers to current vue instance a IssueChart component
// having the following profile
//  <issue-chart :issuestat="some issuestat"></issue-chart>
func RegisterIssueChartComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewIssueChartComp

	o.AddProp("issuestat")

	o.Template = `
	<div>
		<h5 v-if="issuestat.href">Evolution des temps : <a :href="issuestat.href" target="_blank">{{issuestat.issue}}</a></h5>
		<h4 v-else>Evolution des temps : {{issuestat.issue}}</h4>
		<div class="issuechart" ref="container" :style="SetStyle()"></div>
	</div>`

	o.AddMethodWithReturnValue("SetStyle", func(vm *vue.ViewModel, args []*js.Object) interface{} {
		ic := &IssueChartComp{Object: vm.Object}
		if ic.IssueStat.HRef == "" {
			return "width:100%; height:200px;"
		}
		return "width:100%; height:150px;"
	})

	o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		ic := &IssueChartComp{Object: vm.Object}
		is := ic.IssueStat

		startDate := dateFromJSString(is.StartDate)

		chartdesc := js.M{
			"chart": js.M{
				"backgroundColor": "#F7F7F7",
				"type":            "line",
			},
			"title": js.M{
				"text": nil,
			},
			//"xAxis": js.M{
			//	"categories": is.Dates,
			//	"tickPixelInterval" : 400,
			//},
			"xAxis": js.M{
				"type": "datetime",
				"dateTimeLabelFormats": js.M{
					"day": "%e %b",
				},
			},
			"yAxis": js.M{
				"title": js.M{
					"text": "Jours",
				},
			},
			"legend": js.M{
				"layout":        "vertical",
				"align":         "right",
				"verticalAlign": "top",
			},
			"plotOptions": js.M{
				"series": js.M{
					"allowPointSelect": false,
					"pointStart":       startDate,
					"pointInterval":    24 * 3600 * 1000, // one day
				},
			},
			"series": js.S{
				js.M{
					"name": "Passé",
					"data": is.TimeSpent,
				},
				js.M{
					"name": "Restant",
					"data": is.TimeRemaining,
				},
				js.M{
					"name": "Estimé",
					"data": is.TimeEstimated,
				},
			},
		}
		js.Global.Get("Highcharts").Call("chart", vm.Refs.Get("container"), chartdesc)
	})

	return o.NewComponent().Register("issue-chart")
}
