package frontmodel

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

type ProjectStatNames struct {
	*js.Object
	Clients  []string `json:"clients"   js:"clients"`
	Projects []string `json:"projects"  js:"projects"`
}

func NewProjectStatNameFromJS(o *js.Object) *ProjectStatNames {
	psn := &ProjectStatNames{Object: o}
	return psn
}

func (psn *ProjectStatNames) ToClientNameList() []*ValText {
	res := []*ValText{}
	for i, c := range psn.Clients {
		res = append(res, NewValText(c, psn.Projects[i]))
	}
	return res
}

func NewProjectStatName() *ProjectStatNames {
	psn := &ProjectStatNames{Object: js.Global.Get("Object").New()}
	psn.Clients = []string{}
	psn.Projects = []string{}
	return psn
}

func NewProjectStatNameFromList(list []string, sep string) *ProjectStatNames {
	psn := &ProjectStatNames{}
	psn.Clients = make([]string, len(list))
	psn.Projects = make([]string, len(list))
	for i, s := range list {
		elems := strings.Split(s, sep)
		if len(elems) != 2 {
			continue
		}
		psn.Clients[i] = elems[0]
		psn.Projects[i] = elems[1]
	}
	return psn
}
