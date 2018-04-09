package frontmodel

import "github.com/gopherjs/gopherjs/js"

type Portfolio struct {
	*js.Object
	Projects []Project `js:"projects"`
}

func NewPortfolioForBE() *Portfolio {
	return &Portfolio{Projects: []Project{}}
}

func NewPortfolio() *Portfolio {
	return &Portfolio{
		Object:   js.Global.Get("Object").New(),
		Projects: []Project{},
	}
}
