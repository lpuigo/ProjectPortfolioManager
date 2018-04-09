package frontmodel

import (
	"github.com/gopherjs/gopherjs/js"
	"sort"
	"strings"
)

// TextFiltered returns true if prj is not hiden by given filter
func TextFiltered(prj *Project, filter string) bool {
	expected := true
	if filter == "" {
		return true
	}
	if strings.HasPrefix(filter, `\`) {
		if len(filter) > 1 { // prevent from filtering all when only '\' is entered
			expected = false
		}
		filter = filter[1:]
	}
	return prj.Contains(filter) == expected
}

type FilterItem struct {
	*js.Object
	Text     string `js:"text"`
	Count    int    `js:"count"`
	Selected bool   `js:"selected"`
}

func NewFilterItem(text string) *FilterItem {
	fi := &FilterItem{Object: js.Global.Get("Object").New()}
	fi.Text = text
	fi.Count = 0
	fi.Selected = true
	return fi
}

type ColFilter struct {
	*js.Object
	Filters   []*FilterItem           `js:"filters"`
	Index     map[string]*FilterItem  `js:"index"`
	ColName   string                  `js:"colname"`
	GetAttrib func(p *Project) string `js:"getAttrib"`
}

func NewColFilter(colName string) *ColFilter {
	fcl := &ColFilter{Object: js.Global.Get("Object").New()}
	fcl.Filters = []*FilterItem{}
	fcl.Index = map[string]*FilterItem{}
	fcl.ColName = colName
	fcl.GetAttrib = func(p *Project) string {
		panic("Implement me")
		return ""
	}
	return fcl
}

func (cf *ColFilter) AddFilter(f *FilterItem) {
	cf.Filters = append(cf.Filters, f)
	cf.Object.Get("index").Set(f.Text, f)
}

func (cf *ColFilter) Sort() {
	cf.Object.Get("filters").Call("sort", func(a, b *js.Object) int {
		fa, fb := &FilterItem{Object: a}, &FilterItem{Object: b}
		if fa.Text < fb.Text {
			return -1
		} else if fa.Text > fb.Text {
			return 1
		}
		return 0
	})
}

func (cf *ColFilter) ResetFilterList() {
	cf.Filters = []*FilterItem{}
}

func (cf *ColFilter) ResetFilterCounters() {
	for _, fi := range cf.Filters {
		fi.Count = 0
	}
}

type ColFilterGroup struct {
	*js.Object
	CFs []*ColFilter `js:"cfls"`
	//Refresh bool         `js:"refresh"`
}

func NewColFilterGroup() *ColFilterGroup {
	cfg := &ColFilterGroup{Object: js.Global.Get("Object").New()}
	cfg.CFs = []*ColFilter{}
	//cfg.Refresh = false
	return cfg
}

func (cfg *ColFilterGroup) AddColFilter(colName string, f func(p *Project) string) {
	cflStatus := NewColFilter(colName)
	cflStatus.GetAttrib = f
	cfg.CFs = append(cfg.CFs, cflStatus)
}

// GetColFilter returns ColFilter related to colName.
// Panics if colName is not found
func (cfg *ColFilterGroup) GetColFilter(colName string) *ColFilter {
	for _, cf := range cfg.CFs {
		if cf.ColName == colName {
			return cf
		}
	}
	panic("Can not find ColFilter related to " + colName)
	return nil
}

func (cfg *ColFilterGroup) RefreshFilterCounters(prjs []*Project) {
	// reset counters
	for _, cf := range cfg.CFs {
		cf.ResetFilterCounters()
	}
	// parse projects and update counters
	for _, p := range prjs {
		for _, cf := range cfg.CFs {
			cf.Index[cf.GetAttrib(p)].Count++
		}
	}
}

// InitCFL populates ColFilter related to colName
func (cfg *ColFilterGroup) InitCFL(prjs []*Project, colName string) []*FilterItem {
	cf := cfg.GetColFilter(colName)
	if len(cf.Filters) > 0 {
		return cf.Filters
	}

	// parse prjs and create
	count := map[string]int{}
	attribs := []string{}
	for _, p := range prjs {
		attrib := cf.GetAttrib(p)
		if _, exist := count[attrib]; !exist {
			attribs = append(attribs, attrib)
		}
		count[attrib]++
	}
	sort.Strings(attribs)
	for _, attrib := range attribs {
		filterItem := NewFilterItem(attrib)
		filterItem.Count = count[attrib]
		filterItem.Selected = true
		cf.AddFilter(filterItem)
		//cf.Filters = append(cf.Filters, filterItem)
	}
	//cfg.Refresh = true
	return cf.Filters
}

func (cfg *ColFilterGroup) SelectAll(colName string) {
	cf := cfg.GetColFilter(colName)
	for _, fi := range cf.Filters {
		fi.Selected = true
	}
	//cfg.Refresh = true
}

func (cfg *ColFilterGroup) SelectOnly(colName string, text string) {
	cf := cfg.GetColFilter(colName)
	for _, fi := range cf.Filters {
		fi.Selected = (fi.Text == text)
	}
}

func (cfg *ColFilterGroup) InvertSelect(colName string, text string) {
	cf := cfg.GetColFilter(colName)
	if fi, found := cf.Index[text]; found {
		fi.Selected = !fi.Selected
	}
}

// ColFiltered returns true if prj is not hiden by given ColFilterGroup
func (cfg *ColFilterGroup) ColFiltered(prj *Project) bool {
	for _, cf := range cfg.CFs {
		fAttrib := cf.GetAttrib(prj)
		if fi, found := cf.Index[fAttrib]; found && !fi.Selected {
			return false
		}
	}
	return true
}
