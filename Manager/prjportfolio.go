package Manager

import (
	"encoding/json"
	"github.com/lpuig/Novagile/Model"
	"github.com/lpuig/Novagile/Model/IdIndex"
	"os"
)

type PrjPortfolio struct {
	index    IdIndex.IdIndex
	Projects []*Model.Project `json:"projects"`
}

func NewPrjPortfolio() *PrjPortfolio {
	pptf := &PrjPortfolio{Projects: []*Model.Project{}}
	pptf.refreshIndex()
	return pptf
}

// GetPrjById returns the Project with given Id (nil if not found)
func (pptf *PrjPortfolio) GetPrjById(id int) *Model.Project {
	e, found := pptf.index.ById(id)
	if !found {
		return nil
	}
	return pptf.Projects[e]
	//for _, p := range pptf.Projects {
	//	if p.Id == id {
	//		return p
	//	}
	//}
	//return nil
}

func (pptf *PrjPortfolio) refreshIndex() {
	pptf.index = IdIndex.New()
	for pos, p := range pptf.Projects {
		pptf.index.AddElem(p.Id, pos)
	}
}

func (pptf *PrjPortfolio) nextId() int {
	nid := 1
	if len(pptf.Projects) == 0 {
		return nid
	}
	for _, p := range pptf.Projects {
		if nid <= p.Id {
			nid = p.Id
		}
	}
	return nid + 1
}

// AddPrj adds the given project to the PrjPortfolio (new project p is assigned a new Id, but no unicity control is being undertaken for PrjPortfolio consistency)
func (pptf *PrjPortfolio) AddPrj(p *Model.Project) *Model.Project {
	p.Id = pptf.nextId()
	pptf.index.AddElem(p.Id, len(pptf.Projects))
	pptf.Projects = append(pptf.Projects, p)
	return p
}

// DeletePrj deletes the given project id from the PrjPortfolio
//
// If pId if not found, it's a no-op and DeletePrj returns false
func (pptf *PrjPortfolio) DeletePrj(pId int) bool {
	pos, found := pptf.index.ById(pId)
	if !found {
		return false
	}
	pptf.Projects = append(pptf.Projects[:pos], pptf.Projects[pos+1:]...)
	pptf.refreshIndex()
	return true
	//for i, p := range pptf.Projects {
	//	if p.Id == pId {
	//		pptf.Projects = append(pptf.Projects[:i], pptf.Projects[i+1:]...)
	//		return true
	//	}
	//}
	//return false
}

func (pptf *PrjPortfolio) GetPrjClientName(sep string) map[string]bool {
	res := make(map[string]bool)
	for _, p := range pptf.Projects {
		res[sep+p.Client+sep+p.Name] = true
	}
	return res
}

func (pptf *PrjPortfolio) String() string {
	res := "PrjPortfolio :\n"
	for _, p := range pptf.Projects {
		res += p.String()
	}
	return res
}

func NewPrjPortfolioFromJSONFile(f string) (*PrjPortfolio, error) {
	var ptf PrjPortfolio
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&ptf)
	if err != nil {
		return nil, err
	}
	ptf.refreshIndex()
	return &ptf, nil
}

// WriteJsonFile persists ptf (JSON  Marshal) to file with path f
func (pptf *PrjPortfolio) WriteJsonFile(f string) error {
	w, err := os.Create(f)
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(pptf)
}
