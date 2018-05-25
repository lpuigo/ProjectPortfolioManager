package projecttree

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/novagile/src/client/tools"
)

const template string = `
<el-tree
    :data="nodes"
    :props="nodeProps"
	v-loading="nodes.length == 0"
>
    <span class="custom-tree-node" slot-scope="{ node, data }">
		<el-tooltip v-if="data.href" :content="data.summary" placement="right" effect="light">
			<span><a :href="data.href" target="_blank" class="custom-node-name">{{node.label}}</a>
			&nbsp({{data.tothour | FormatFloat(1)}} h)</span>
		</el-tooltip>
		<span v-else class="custom-node-name">{{node.label}}</span>

		<span class="custom-hours-row reduce">
			<span class="hours-table">
				<span class="hours-cell right">{{data.lhour + data.nlhour | FormatFloat(1)}} h</span>
				<span class="hours-cell"></span>
				<el-progress 
						v-if="data.level==0"	
						class="hours-cell large"
						:text-inside="true"
						:stroke-width="18"
						:percentage="data.lhour *100 / (data.lhour + data.nlhour) | FormatFloat"
				></el-progress>
				<el-progress 
						v-else
						class="hours-cell large"
						color="#4abbbd"
						:text-inside="true"
						:stroke-width="15"
						:percentage="(data.lhour + data.nlhour)*100 / (data.parent.lhour + data.parent.nlhour) | FormatFloat"
				></el-progress>
			</span>
		</span>
    </span>
</el-tree>
`

func Register() {
	hvue.NewComponent("project-tree",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Props("nodes"),
		hvue.Template(template),
		//hvue.Component("hours-row", hoursrow.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewHoursTreeCompModel(vm)
		}),
		hvue.MethodsOf(&ProjectTreeCompModel{}),
		hvue.Filter("FormatFloat", func(vm *hvue.VM, value *js.Object, args ...*js.Object) interface{} {
			h := value.Float()
			prec := 0
			if len(args) > 0 {
				prec = args[0].Int()
			}
			return strconv.FormatFloat(h, 'f', prec, 64)
		}),
	}
}

type ProjectTreeCompModel struct {
	*js.Object

	Nodes     []*Node `js:"nodes"`
	NodeProps js.M    `js:"nodeProps"`

	VM *hvue.VM `js:"VM"`
}

func NewHoursTreeCompModel(vm *hvue.VM) *ProjectTreeCompModel {
	htcm := &ProjectTreeCompModel{Object: tools.O()}
	htcm.NodeProps = js.M{
		"children": "children",
		"label":    "label",
	}

	htcm.Nodes = []*Node{}
	htcm.VM = vm
	return htcm
}
