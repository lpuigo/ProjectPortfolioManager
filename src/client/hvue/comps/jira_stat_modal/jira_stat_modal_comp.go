package jira_stat_modal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	jsn "github.com/lpuig/element/model/jirastatnode"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/hourstree"
	"github.com/lpuig/novagile/src/client/hvue/comps/jira_stat_modal/node"
	"github.com/lpuig/novagile/src/client/tools"
)

func Register() {
	hvue.NewComponent("jira-stat-modal",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Component("hours-tree", hourstree.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewJiraStatModalModel(vm)
		}),
		hvue.MethodsOf(&JiraStatModalModel{}),
	}
}

type JiraStatModalModel struct {
	*js.Object

	Visible bool     `js:"visible"`
	VM      *hvue.VM `js:"VM"`

	Nodes []*node.HoursNode `js:"nodes"`
}

func NewJiraStatModalModel(vm *hvue.VM) *JiraStatModalModel {
	jsmm := &JiraStatModalModel{Object: tools.O()}
	jsmm.Visible = false
	jsmm.VM = vm

	jsmm.Nodes = []*node.HoursNode{}
	return jsmm
}

func (jsmm *JiraStatModalModel) Show() {
	go jsmm.GetNodes()
	jsmm.Visible = true
}

func (jsmm *JiraStatModalModel) Hide() {
	jsmm.Visible = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Button Methods

func (jsmm *JiraStatModalModel) HandleNodeClick() {

}

func (jsmm *JiraStatModalModel) GetNodes() {
	jsmm.Nodes = []*node.HoursNode{}
	jsns := callGetJiraStat()
	res := []*node.HoursNode{}

	var teamnode *node.HoursNode
	team := ""

	jsns.Call("forEach", func(jsn *jsn.JiraStatNode) {
		if jsn.Team != team {
			team = jsn.Team
			teamnode = node.NewHoursNode(team, nil, 0)
			res = append(res, teamnode)
		}
		teamnode.AddChild(node.NewHoursNode(jsn.Author, jsn.HourLogs, 40))
	})

	jsmm.Nodes = res
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

//////////////////////////////////////////////////////////////////////////////////////////////
// Get Client-Name List
var jiradata string = `[{"team":"Novagile DEV","author":"a.juffet","hour_logs":[0,0,0,0,0,0,12,0,0,0,0,0,0,2.5]},{"team":"Novagile DEV","author":"c.hardouin","hour_logs":[0,0,0,0,16,40,40,40,40,40,40,40,40,16]},{"team":"Novagile DEV","author":"d.zhong","hour_logs":[0,0,0,0,0,0,0,4,0,0,0,0,0,0]},{"team":"Novagile DEV","author":"f.couffe","hour_logs":[32,40,40,40,40,40,40,40,40,40,40,40,40,8]},{"team":"Novagile DEV","author":"g.sebastiani","hour_logs":[0.75,14,35.75,13.5,26.75,26.5,25.5,41,4.8332999999999995,2.5,30,37.4999,36.9167,32]},{"team":"Novagile DEV","author":"h.atig","hour_logs":[0,16,8,0,10,1,32,0,32,28,24,0,42,7]},{"team":"Novagile DEV","author":"j.georget","hour_logs":[18.8333,16.6667,14.6167,5.6667,0,9.166699999999999,8.2,2.5668,9.3168,3.1333,2.6667,11.5,31.2667,3]},{"team":"Novagile DEV","author":"l.liu","hour_logs":[0,0,0,0,0,0,0,0,0,0,8,35,36,8]},{"team":"Novagile DEV","author":"l.manns","hour_logs":[28,40,38,40,35,40,41,46,39,40,41,40,40,10]},{"team":"Novagile DEV","author":"s.nay","hour_logs":[8.749999999999998,4.5,32,15,32,40,40,40,40,0,30,35.5,29.5,3]},{"team":"Novagile DEV","author":"s.parent","hour_logs":[32.75,42,40,45,38.75,57.25,41,40.25,40.25,40,42.75,43.75,39.75,17.75]},{"team":"Novagile DEV","author":"t.capon","hour_logs":[28.000099999999996,35.1667,35.3834,35.1999,34.9499,35.666599999999995,37.08369999999999,36.5833,0,35.3334,37.050000000000004,36,35.5,0]},{"team":"Novagile DEV","author":"t.planchon","hour_logs":[28.7,39.5836,31.2334,40,40,36.5501,33.4166,36.633300000000006,29.9664,35.6834,31.583399999999997,36.40019999999999,35.8833,5.450099999999998]},{"team":"Novagile PMO","author":"l.puig","hour_logs":[28,35,38,37,39,39.5,40,40,40,40,40,40,40,16]},{"team":"Novagile PMO","author":"m.tabuy","hour_logs":[2,15,9,17,5.5,10,7.5,9,7,0,0,0,40,16]},{"team":"Novagile PS France","author":"l.fadlane","hour_logs":[32,32,40,40,40,40,34,32,35,39.5,37,0,0,0]},{"team":"Novagile PS France","author":"v.vanelsuve","hour_logs":[0,0,0,0,0,2,0,0,0,0,28,0,0,0]},{"team":"Novagile PS France","author":"y.maigrez","hour_logs":[0,0,0,1,0,3,6,0,0,0,0,0,0,0]},{"team":"france-dsi-rd","author":"a.quessada","hour_logs":[0,0,0,40,40,40,40,40,40,33,0,3,0,0]},{"team":"france-dsi-rd","author":"e.giner","hour_logs":[32,40,40,40,40,40,40,40,27,23,8,22,26,12]},{"team":"france-dsi-rd","author":"f.vaucelle","hour_logs":[32,40,40,40,40,40,40,40,40,8,8,3,16,0]},{"team":"france-dsi-rd","author":"s.akkaoui","hour_logs":[8,40,0,8,16,0,0,0,8,8,8,3,17,0]},{"team":"france-dsi-rd","author":"v.gabou","hour_logs":[32,40,40,40,40,40,40,40,32,0,16,19,40,32]}]
`

func callGetJiraStat() *js.Object {
	// TODO replace with XHR query
	return js.Global.Get("JSON").Call("parse", jiradata)
}
