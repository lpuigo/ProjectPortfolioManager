package Comps

import (
	"github.com/gopherjs/gopherjs/js"
	//fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/oskca/gopherjs-vue"
	//charts "github.com/cnguy/gopherjs-frappe-charts"
)

type IssueChartComp struct {
	*js.Object
	//PrjStat *fm.ProjectStat `js:"projectstat"`
}

func NewIssueChartComp() *IssueChartComp {
	ic := &IssueChartComp{Object: js.Global.Get("Object").New()}
	//ic.PrjStat = nil
	return ic
}

func (ic *IssueChartComp) SetupChart(this *js.Object) {
	//chartData := charts.NewChartData()
	//chartData.Labels = ic.PrjStat.Dates
	//chartData.Datasets = make([]*charts.Dataset, len(ic.PrjStat.Issues)*3)
	//for num, issueName := range ic.PrjStat.Issues {
	//	chartData.Datasets[num+0] = charts.NewDataset(issueName+" Spent", ic.PrjStat.TimeSpent[num])
	//	chartData.Datasets[num+1] = charts.NewDataset(issueName+" Remaining", ic.PrjStat.TimeRemaining[num])
	//	chartData.Datasets[num+2] = charts.NewDataset(issueName+" Estimated", ic.PrjStat.TimeEstimated[num])
	//}
	//
	//_ = charts.NewLineChart(this.String(), chartData).WithHeight(300).SetShowDots(true).Render()
}

// RegisterIssueChartComp registers to current vue intance a IssueChart component
// having the following profile
//  <issue-chart :projectstat="some projectstat"></issue-chart>
func RegisterIssueChartComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewIssueChartComp

	o.AddProp("projectstat")

	o.Template = `
	<div><pre>Stats : {{projectstat}}</pre></div>`

	o.AddMethod("Render", func(vm *vue.ViewModel, args []*js.Object) {
		println("vm.Object", vm.Object)
		//ic := &IssueChartComp{Object: vm.Object}
		//ic.SetupChart(vm.El)
	})

	//o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
	//	ic := &IssueChartComp{Object:vm.Object}
	//	ic.Render(vm.El.String())
	//})

	return o.NewComponent().Register("issue-chart")
}
