package bars_chart

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/prjptf/src/client/tools"
)

func Register() {
	hvue.NewComponent("bars-chart",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("infos"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewBarsChartModel(vm)
		}),
		hvue.MethodsOf(&BarsChartModel{}),
		hvue.Mounted(func(vm *hvue.VM) {
			bcm := &BarsChartModel{Object: vm.Object}
			bcm.showChart()
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

type Infos struct {
	*js.Object
	Weeks  []string `js:"weeks"`
	Series []*Serie `js:"series"`
}

func NewInfos(weeks []string, series []*Serie) *Infos {
	i := &Infos{Object: tools.O()}
	i.Weeks = weeks
	i.Series = series
	return i
}

type BarsChartModel struct {
	*js.Object
	VM    *hvue.VM `js:"VM"`
	Infos *Infos   `js:"infos"`
}

func NewBarsChartModel(vm *hvue.VM) *BarsChartModel {
	bcm := &BarsChartModel{Object: tools.O()}
	bcm.VM = vm

	bcm.Infos = nil
	return bcm
}

func (bcm *BarsChartModel) SetStyle() string {
	return "width:100%; height:550px;"
}

func (bcm *BarsChartModel) Refresh(infos *Infos) {
	bcm.Infos = infos
	bcm.showChart()
}

func (bcm *BarsChartModel) showChart() {
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
			"categories": bcm.Infos.Weeks,
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
		"series": bcm.Infos.Series,
	}
	js.Global.Get("Highcharts").Call("chart", bcm.VM.Refs("container"), chartdesc)
}
