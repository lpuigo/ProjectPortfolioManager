package comps

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	"github.com/oskca/gopherjs-vue"
)

type ColTitleWithFilterComp struct {
	ColTitleComp
	ColFilterGrp *fm.ColFilterGroup `js:"colfilter"`
	Projects     []*fm.Project      `js:"projects"`
}

func NewColTitleWithFilterComp() *ColTitleWithFilterComp {
	ct := &ColTitleWithFilterComp{ColTitleComp: *NewColTitleComp()}
	ct.ColFilterGrp = nil
	ct.Projects = nil
	return ct
}

func NewColTitleWithFilterCompFromJS(o *js.Object) *ColTitleWithFilterComp {
	return &ColTitleWithFilterComp{ColTitleComp: ColTitleComp{Object: o}}
}

func (ct *ColTitleWithFilterComp) RefreshCounters() bool {
	go ct.ColFilterGrp.RefreshFilterCounters(ct.Projects)
	return true
}

// RegisterColTitleComp registers to current vue intance a EditProjectModal component
// having the following profile
//  <th is="col-title-withfilter"
//     :sortlist.sync="some_[]*SortCol"
//     :title="string"
//     :colfilter.sync="some_*ColFilterGroup"
//     :projects="some_[]project"
//     :iclass="icon class">
//  </th>
func RegisterColTitleWithFilterComp() *vue.Component {
	var jq = jquery.NewJQuery
	o := vue.NewOption()
	o.Data = NewColTitleWithFilterComp

	o.AddProp("sortlist", "title", "colfilter", "projects", "iclass")

	o.Template = `
	<th>
		<div ref="coltitleDD" class="ui dropdown item">
			<i v-if="iclass" :class="iclass"></i>
			<span>{{title}}</span>
			<!-- TODO filter icon if filter is activated -->
			<span v-if="sortedby"><i :class="sortedby"></i>{{rank+1}}</span>
			<div class="menu">
				<div class="item" @click="sortthis($event)">Trier</div>
				<div class="divider"></div>
				<div class="item" @click="selectAll">Tout sélectionner</div>
				<div v-for="(fi, index) in filterItemList" :key="index" class="item" :data-value="fi.text" @click="changeSelection(fi, $event)">
					<i v-if="fi.selected" class="checkmark box icon"></i>
					<i v-else class="square outline icon"></i>
					<div class="ui mini right pointing label">{{fi.count}}</div>
					{{fi.text}}
				</div>
			</div>
		</div>
	</th>`

	o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		ct := NewColTitleWithFilterCompFromJS(vm.Object)
		jq(vm.Refs.Get("coltitleDD")).
			Call("dropdown", js.M{
				"on":     "hover",
				"action": "nothing",
				"onShow": ct.RefreshCounters,
			})
	})

	o.AddComputed("sortedby", func(vm *vue.ViewModel) interface{} {
		ct := NewColTitleWithFilterCompFromJS(vm.Object)
		return ct.sortedBy()
	})

	o.AddComputed("filterItemList", func(vm *vue.ViewModel) interface{} {
		ct := NewColTitleWithFilterCompFromJS(vm.Object)
		filterItems := ct.ColFilterGrp.InitCFL(ct.Projects, vm.Get("title").String())
		//vm.Emit("update:colfilter", ct.ColFilterGrp)
		return filterItems
	})

	o.AddMethod("sortthis", func(vm *vue.ViewModel, args []*js.Object) {
		ct := NewColTitleWithFilterCompFromJS(vm.Object)
		ct.switchSortedBy(args[0].Get("ctrlKey").Bool())
		vm.Emit("update:sortlist", ct.SortList)
	})

	o.AddMethod("selectAll", func(vm *vue.ViewModel, args []*js.Object) {
		ct := NewColTitleWithFilterCompFromJS(vm.Object)
		ct.ColFilterGrp.SelectAll(vm.Get("title").String())
		vm.Emit("update:colfilter", ct.ColFilterGrp)
	})

	o.AddMethod("changeSelection", func(vm *vue.ViewModel, args []*js.Object) {
		ct := NewColTitleWithFilterCompFromJS(vm.Object)
		fi := &fm.FilterItem{Object: args[0]}
		if args[1].Get("ctrlKey").Bool() {
			ct.ColFilterGrp.SelectOnly(vm.Get("title").String(), fi.Text)
		} else {
			ct.ColFilterGrp.InvertSelect(vm.Get("title").String(), fi.Text)
		}
		//fi.Selected = !fi.Selected
		//ct.ColFilterGrp.Refresh = true
		vm.Emit("update:colfilter", ct.ColFilterGrp)
	})

	return o.NewComponent().Register("col-title-withfilter")
}
