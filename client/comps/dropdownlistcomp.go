package comps

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	fm "github.com/lpuig/novagile/client/frontmodel"
	"github.com/oskca/gopherjs-vue"
)

type DropDownListComp struct {
	*js.Object
	Selected   string        `js:"selected"`
	ListValues []*fm.ValText `js:"listvalues"`
}

func NewDropDownListComp() *DropDownListComp {
	a := &DropDownListComp{Object: js.Global.Get("Object").New()}
	a.Selected = ""
	return a
}

const (
	TemplateDropDownListComp = `
	<select class="ui search selection dropdown" v-model="selected">
		<option value="">{{defaulttext}}</option>
		<option v-for="vt in listvalues" :key="vt.value" :value="vt.value">{{vt.text}}</option>
	</select>
`
)

// RegisterDropDownListComp registers to current vue intance a DropDownListComp component
// having the following profile
//  <dropdown-list :listvalues="some_[]*ValText" defaulttext="text" :selected.sync="binded_variable"></dropdown-list>
func RegisterDropDownListComp() *vue.Component {
	var jq = jquery.NewJQuery

	o := vue.NewOption()
	o.Template = TemplateDropDownListComp
	o.Data = NewDropDownListComp

	o.AddProp("selected", "listvalues", "defaulttext")

	o.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		ddl := &DropDownListComp{Object: vm.Object}
		jq(vm.El).Call("dropdown", js.M{"onChange": func(v, t, s *js.Object) {
			vm.Emit("update:selected", v)
		}})
		jq(vm.El).Call("dropdown", "set selected", ddl.Selected)
	})

	//TODO change this with a watcher so selection update is automatic when selected change
	o.AddMethod("changeSelected", func(vm *vue.ViewModel, args []*js.Object) {
		ddl := &DropDownListComp{Object: vm.Object}
		val := args[0].String()
		if fm.IsInValTextList(val, ddl.ListValues) {
			jq(vm.El).Call("dropdown", "set selected", args[0])
		} else {
			jq(vm.El).Call("dropdown", "clear")
		}
	})

	c := o.NewComponent()

	return c.Register("dropdown-list")
}
