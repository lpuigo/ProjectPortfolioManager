package sre_chart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/lpuig/novagile/src/client/hvue/tools"
	"time"
)

func Register() {
	hvue.NewComponent("sre-chart",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("issuestat"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewSREChartModel(vm)
		}),
		hvue.MethodsOf(&SREChartModel{}),
		hvue.Mounted(func(vm *hvue.VM) {
			setChart(vm)
		}),
	}
}

type SREChartModel struct {
	*js.Object

	IssueStat *fm.IssueStat `js:"issuestat"`
	VM        *hvue.VM      `js:"VM"`
}

func NewSREChartModel(vm *hvue.VM) *SREChartModel {
	scm := &SREChartModel{Object: tools.O()}

	scm.VM = vm
	return scm
}

func (scm *SREChartModel) SetStyle() string {
	return "width:100%; height:250px;"
}

func setChart(vm *hvue.VM) {
	ic := &SREChartModel{Object: vm.Object}
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
				"text": "Days",
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
				"marker":           js.M{"enabled": false},
				"animation":        false,
			},
		},
		"series": js.S{
			js.M{
				"name":      "Spent",
				"color":     "#51A825",
				"lineWidth": 5,
				"data":      is.TimeSpent,
			},
			js.M{
				"name":  "Remaining",
				"color": "#EC8E00",
				"data":  is.TimeRemaining,
			},
			js.M{
				"name":  "Estimated",
				"color": "#89AFD7",
				"data":  is.TimeEstimated,
			},
		},
	}
	js.Global.Get("Highcharts").Call("chart", vm.Refs("container"), chartdesc)
}

const TimeJSLayout string = "2006-01-02"

func dateFromJSString(s string) int64 {
	res := time.Time{}
	res, _ = time.Parse(TimeJSLayout, s)
	return res.Unix() * 1000
}
