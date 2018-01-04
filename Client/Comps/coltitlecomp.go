package Comps

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/oskca/gopherjs-vue"
)

type ColTitleComp struct {
	*js.Object
	SortList []*fm.SortCol `js:"sortlist"`
	Title    string        `js:"title"`
	SortRank int           `js:"rank"`
}

func NewColTitleComp() *ColTitleComp {
	ct := &ColTitleComp{Object: js.Global.Get("Object").New()}
	ct.Title = ""
	ct.SortList = nil
	ct.SortRank = -1
	return ct
}

func (ct *ColTitleComp) sortedBy() string {
	if len(ct.SortList) == 0 {
		return ""
	}
	if ct.setSortRank() < 0 {
		return ""
	}
	if ct.SortList[ct.SortRank].Asc {
		return "sort content ascending icon"
	}
	return "sort content descending icon"
}

func (ct *ColTitleComp) setSortRank() int {
	rank := -1
	for i, sc := range ct.SortList {
		if sc.Name == ct.Title {
			rank = i
			break
		}
	}
	ct.SortRank = rank
	return rank
}

// switchSortedBy updates ct.SortList with the given column name.
//
// if addmode is true, current col name is added (or reversed) in the current sort list,
// otherwise, current col name replaces existing SortList
func (ct *ColTitleComp) switchSortedBy(addmode bool) {
	defer ct.setSortRank()
	if addmode {
		if ct.SortRank < 0 {
			ct.SortList = append(ct.SortList, fm.NewSortCol(ct.Title, true))
			return
		}
		ct.SortList[ct.SortRank].Asc = !ct.SortList[ct.SortRank].Asc
		return
	}
	// replace mode
	if ct.SortRank >= 0 {
		ct.SortList = []*fm.SortCol{fm.NewSortCol(ct.Title, !ct.SortList[ct.SortRank].Asc)}
		return
	}
	ct.SortList = []*fm.SortCol{fm.NewSortCol(ct.Title, true)}
	return
}

// RegisterColTitleComp registers to current vue intance a EditProjectModal component
// having the following profile
//  <col-title :sortlist.sync="some_[]*SortCol" :name="string" :isdate="bool" :iclass="icon class"></col-title>
func RegisterColTitleComp() *vue.Component {
	o := vue.NewOption()
	o.Data = NewColTitleComp

	o.AddProp("sortlist", "title", "iclass")

	o.Template = `
	<th @click="sortThis($event)">
		<i v-if="iclass" :class="iclass"></i>
		<span>{{title}}</span>
		<span v-if="sortedby"><i :class="sortedby"></i>{{rank+1}}</span>
	</th>`

	o.AddComputed("sortedby", func(vm *vue.ViewModel) interface{} {
		ct := &ColTitleComp{Object: vm.Object}
		return ct.sortedBy()
	})

	o.AddMethod("sortThis", func(vm *vue.ViewModel, args []*js.Object) {
		ct := &ColTitleComp{Object: vm.Object}
		ct.switchSortedBy(args[0].Get("ctrlKey").Bool())
		vm.Emit("update:sortlist", ct.SortList)
	})

	return o.NewComponent().Register("col-title")
}
