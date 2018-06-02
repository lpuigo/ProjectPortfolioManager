package hourstree

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/prjptf/src/client/hvue/comps/jira_stat_modal/hoursrow"
	"github.com/lpuig/prjptf/src/client/tools"
)

const template = `
<el-tree
    :data="nodes"
    :props="nodeProps"
	v-loading="nodes.length == 0"
>
    <span class="custom-tree-node" slot-scope="{ node, data }">
        <span class="custom-node-name">{{ node.label }}</span>
        <hours-row :hours="data.hours" :hmax="data.maxHour"></hours-row>
    </span>
</el-tree>
`

func Register() {
	hvue.NewComponent("hours-tree",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Props("nodes"),
		hvue.Template(template),
		hvue.Component("hours-row", hoursrow.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewHoursTreeCompModel(vm)
		}),
		hvue.MethodsOf(&HoursTreeCompModel{}),
	}
}

type HoursTreeCompModel struct {
	*js.Object

	Nodes     []*Node `js:"nodes"`
	NodeProps js.M    `js:"nodeProps"`

	VM *hvue.VM `js:"VM"`
}

func NewHoursTreeCompModel(vm *hvue.VM) *HoursTreeCompModel {
	htcm := &HoursTreeCompModel{Object: tools.O()}
	htcm.NodeProps = js.M{
		"children": "children",
		"label":    "label",
	}

	htcm.Nodes = []*Node{}
	htcm.VM = vm
	return htcm
}
