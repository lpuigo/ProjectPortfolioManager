package manager

import (
	"encoding/json"
	"github.com/lpuig/novagile/model"
	"github.com/lpuig/novagile/model/idindex"
	"os"
)

type StatPortfolio struct {
	index idindex.IdIndex
	Stats []*model.ProjectStat `json:"stats"`
}

func NewStatPortfolio() *StatPortfolio {
	sptf := &StatPortfolio{Stats: []*model.ProjectStat{}}
	sptf.refreshIndex()
	return sptf
}

// GetStatById returns the ProjectStat with given Id (nil if not found)
func (sptf *StatPortfolio) GetStatById(id int) *model.ProjectStat {
	e, found := sptf.index.ById(id)
	if !found {
		return nil
	}
	return sptf.Stats[e]
	//for _, p := range sptf.Stats {
	//	if p.Id == id {
	//		return p
	//	}
	//}
	//return nil
}

func (sptf *StatPortfolio) refreshIndex() {
	sptf.index = idindex.New()
	for pos, s := range sptf.Stats {
		sptf.index.AddElem(s.Id, pos)
	}
}

func (sptf *StatPortfolio) nextId() int {
	nid := 1
	if len(sptf.Stats) == 0 {
		return nid
	}
	for _, s := range sptf.Stats {
		if nid <= s.Id {
			nid = s.Id
		}
	}
	return nid + 1
}

// AddProjectStat adds the given ProjectStat to the StatPortfolio
func (sptf *StatPortfolio) AddProjectStat(ps *model.ProjectStat) *model.ProjectStat {
	sptf.index.AddElem(ps.Id, len(sptf.Stats))
	sptf.Stats = append(sptf.Stats, ps)
	return ps
}

// DeleteProjectStat deletes the given ProjectStat id from the StatPortfolio
//
// If pId if not found, it's a no-op and DeleteProjectStat returns false
func (sptf *StatPortfolio) DeleteProjectStat(pId int) bool {
	pos, found := sptf.index.ById(pId)
	if !found {
		return false
	}
	sptf.Stats = append(sptf.Stats[:pos], sptf.Stats[pos+1:]...)
	sptf.refreshIndex()
	return true
	//for i, s := range sptf.Stats {
	//	if s.Id == pId {
	//		sptf.Stats = append(sptf.Stats[:i], sptf.Stats[i+1:]...)
	//		return true
	//	}
	//}
	//return false
}

func (sptf *StatPortfolio) String() string {
	res := "StatPortfolio :\n"
	for _, s := range sptf.Stats {
		res += s.String()
	}
	return res
}

func NewStatPortfolioFromJSONFile(f string) (*StatPortfolio, error) {
	var sptf StatPortfolio
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&sptf)
	if err != nil {
		return nil, err
	}
	sptf.refreshIndex()
	return &sptf, nil
}

// WriteJsonFile persists ptf (JSON  Marshal) to file with path f
func (sptf *StatPortfolio) WriteJsonFile(f string) error {
	w, err := os.Create(f)
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(sptf)
}
