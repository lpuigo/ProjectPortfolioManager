package main

import (
	"github.com/gopherjs/gopherjs/js"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	"github.com/oskca/gopherjs-vue"
	"strconv"
)

// ChargeFilterRegister registers a filter taking Project as argument and returning formated Workload information
func ChargeFilterRegister(name string) {
	cf := vue.NewFilter(func(oldValue *js.Object) (newValue interface{}) {
		p := fm.Project{Object: oldValue}
		if p.ForecastWL > 0 {
			res := strconv.FormatFloat(p.CurrentWL, 'f', 1, 64)
			res += " / "
			res += strconv.FormatFloat(p.ForecastWL, 'f', 1, 64)
			return res
		}
		return "-"
	})
	cf.Register(name)
}

// DateFilterRegister registers a filter taking Date as argument and returning formated date (jj/mm/aaaa)
func DateFilterRegister(name string) {
	cf := vue.NewFilter(func(oldValue *js.Object) (newValue interface{}) {
		return fm.DateString(oldValue.String())
	})
	cf.Register(name)
}
