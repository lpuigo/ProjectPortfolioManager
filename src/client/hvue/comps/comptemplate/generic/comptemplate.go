package generic

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/huckridgesw/hvue"
	"github.com/lpuig/prjptf/src/client/tools"
)

const template = `
`

func Register() {
	hvue.NewComponent("XXX-comp-name",
		ComponentOptions()...,
	)
}

func ComponentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		//hvue.Props("CompPropName"),
		hvue.Template(template),
		//hvue.Component("XXX-inner-comp-name", XXXinnercomp.ComponentOptions()...),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewXXXCompModel(vm)
		}),
		hvue.MethodsOf(&XXXCompModel{}),
	}
}

type XXXCompModel struct {
	*js.Object

	//attribute here

	VM *hvue.VM `js:"VM"`
}

func NewXXXCompModel(vm *hvue.VM) *XXXCompModel {
	htcm := &XXXCompModel{Object: tools.O()}

	//attribute init here

	htcm.VM = vm
	return htcm
}
