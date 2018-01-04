package FrontModel

import "github.com/gopherjs/gopherjs/js"

type ValText struct {
	*js.Object
	Value string `js:"value"`
	Text  string `js:"text"`
}

func NewValText(val, text string) *ValText {
	vt := &ValText{Object: js.Global.Get("Object").New()}
	vt.Value = val
	vt.Text = text
	return vt
}

func IsInValTextList(value string, vtl []*ValText) bool {
	for _, vt := range vtl {
		if vt.Value == value {
			return true
		}
	}
	return false
}
