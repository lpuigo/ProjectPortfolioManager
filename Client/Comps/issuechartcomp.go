package Comps

import (
	"github.com/cnguy/gopherjs-frappe-charts"
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/oskca/gopherjs-vue"
)

type IssueChartComp struct {
	*js.Object
	PrjStat *fm.ProjectStat `js:"projectstat"`
}

func NewIssueChartComp() *IssueChartComp {
	ic := &IssueChartComp{Object: js.Global.Get("Object").New()}
	ic.PrjStat = nil
	return ic
}

func (ic *IssueChartComp) SetupChart(this *js.Object) {
	chartData := charts.NewChartData()
	chartData.Labels = ic.PrjStat.Dates
	chartData.Datasets = []*charts.Dataset{}
	for num, issueName := range ic.PrjStat.Issues {
		chartData.Datasets = append(chartData.Datasets, charts.NewDataset(issueName+" Spent", ic.PrjStat.TimeSpent[num]))
		chartData.Datasets = append(chartData.Datasets, charts.NewDataset(issueName+" Remaining", ic.PrjStat.TimeRemaining[num]))
		chartData.Datasets = append(chartData.Datasets, charts.NewDataset(issueName+" Estimated", ic.PrjStat.TimeEstimated[num]))
	}
	chart := charts.NewBarChart("#chart", chartData).
		WithHeight(300).
		WithColors([]string{"light-blue", "violet", "red"}).
		Render()
	println("chartData", chart)

}

// RegisterIssueChartComp registers to current vue intance a IssueChart component
// having the following profile
//  <issue-chart :projectstat="some projectstat"></issue-chart>
func RegisterIssueChartComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewIssueChartComp

	o.AddProp("projectstat")

	o.Template = `
	<div id="chart" style="background: gray"><pre>Stats : {{projectstat}}</pre></div>`

	o.AddMethod("RenderChart", func(vm *vue.ViewModel, args []*js.Object) {
		ic := &IssueChartComp{Object: vm.Object}
		ic.PrjStat = &fm.ProjectStat{Object: args[0]}
		ic.SetupChart(vm.El)
	})

	//o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
	//	ic := &IssueChartComp{Object:vm.Object}
	//	ic.Render(vm.El.String())
	//})

	return o.NewComponent().Register("issue-chart")
}
