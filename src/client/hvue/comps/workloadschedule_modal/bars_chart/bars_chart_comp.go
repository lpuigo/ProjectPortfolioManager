package bars_chart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	wsr "github.com/lpuig/novagile/src/client/frontmodel/workloadschedulerecord"
	"github.com/lpuig/novagile/src/client/tools"
)

func Register() {
	hvue.NewComponent("bars-chart",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("weeks", "series"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewBarsChartModel(vm)
		}),
		hvue.MethodsOf(&BarsChartModel{}),
		hvue.Mounted(func(vm *hvue.VM) {
			setChart(vm)
		}),
	}
}

type Serie struct {
	*js.Object
	Name  string    `js:"name"`
	Color string    `js:"color"`
	Data  []float64 `js:"data"`
}

func NewSerie(name, color string, data []float64) *Serie {
	s := &Serie{Object: tools.O()}
	s.Name = name
	s.Color = color
	s.Data = data
	return s
}

type BarsChartModel struct {
	*js.Object
	VM       *hvue.VM              `js:"VM"`
	Weeks    []string              `js:"weeks"`
	Series   []*Serie              `js:"series"`
	WrkSched *wsr.WorkloadSchedule `js:"wrkSched"`
}

func NewBarsChartModel(vm *hvue.VM) *BarsChartModel {
	bcm := &BarsChartModel{Object: tools.O()}
	bcm.VM = vm

	bcm.Weeks = nil
	bcm.Series = nil
	return bcm
}

func (scm *BarsChartModel) SetStyle() string {
	return "width:100%; height:550px;"
}

func setChart(vm *hvue.VM) {
	bcm := &BarsChartModel{Object: vm.Object}

	chartdesc := js.M{
		"chart": js.M{
			"backgroundColor": "#F7F7F7",
			"type":            "column",
		},
		"title": js.M{
			"text": nil,
		},
		//"xAxis": js.M{
		//	"categories": is.Dates,
		//	"tickPixelInterval" : 400,
		//},
		"xAxis": js.M{
			"categories": bcm.Weeks,
		},
		"yAxis": js.M{
			"title": js.M{
				"text": "Full-Time Equivalent",
			},
		},
		"legend": js.M{
			"enabled": false,
			//"layout":        "vertical",
			//"align":         "right",
			//"verticalAlign": "top",
		},
		"plotOptions": js.M{
			"column": js.M{
				"stacking": "normal",
				"dataLabels": js.M{
					"enabled": false,
					//"color":   "white",
				},
				"animation": false,
			},
			"series": js.M{
				"pointPadding": -0.25,
			},
		},
		"series": bcm.Series,
	}
	js.Global.Get("Highcharts").Call("chart", vm.Refs("container"), chartdesc)
}
