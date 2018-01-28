package Model

import "sort"

type SituationToDate struct {
	UpdateOn   Date            `json:"update"`
	MileStones map[string]Date `json:"milestones"`
}

func NewSituationToDate() *SituationToDate {
	s := &SituationToDate{}
	s.UpdateOn = Today()
	s.MileStones = make(map[string]Date)

	return s
}

func (s SituationToDate) String() string {
	res := "SituationToDate {\n\tUpdateOn : " + s.UpdateOn.String() + "\n"
	res += "\tMileStones : {"
	keys := []string{}
	for m, _ := range s.MileStones {
		keys = append(keys, string(m))
	}
	sort.Strings(keys)
	for _, m := range keys {
		res += "\n\t\t" + m + " : " + s.MileStones[m].String()
	}
	res += "}\n}\n"
	return res
}

// Clone returns a clone of s (that is a new SituationToDate with same UpdateOn and Situation than s)
func (s *SituationToDate) Clone() *SituationToDate {
	res := NewSituationToDate()
	res.UpdateOn = s.UpdateOn
	for m, d := range s.MileStones {
		res.MileStones[m] = d
	}
	return res
}

// UpdateWith append to s all new or changed MileStones from s2 compared to those in s (s.UpdateOn is also updated)
func (s *SituationToDate) UpdateWith(s2 *SituationToDate) {
	if s2 == nil {
		return
	}
	for m, nd := range s2.MileStones {
		od, found := s.MileStones[m]
		if found {
			if nd.IsZero() {
				delete(s.MileStones, m)
			} else if !od.Equal(nd) {
				// update with changed milestone
				s.MileStones[m] = nd
			}
		} else {
			// update with new milestone
			s.MileStones[m] = nd
		}
	}
	s.UpdateOn = s2.UpdateOn
}

// DifferenceWith returns a New SituationToDate with s2.UpdateOn Date
// and all changed or added MileStones compared to current SituationToDate s
//
// returns nil if s2 is nil or situation unchanged
func (s *SituationToDate) DifferenceWith(s2 *SituationToDate) *SituationToDate {
	if s2 == nil {
		return nil
	}
	std := NewSituationToDate()
	std.UpdateOn = s2.UpdateOn
	// check for new or changed milestone
	for m, nd := range s2.MileStones {
		od, found := s.MileStones[m]
		if !(found && od.Equal(nd)) {
			// Add new or changed milestone
			std.MileStones[m] = nd
		}
	}
	// check for deleted milestone
	for m, _ := range s.MileStones {
		if _, found := s2.MileStones[m]; !found {
			std.MileStones[m] = ZeroDate()
		}
	}
	if len(std.MileStones) == 0 {
		return nil
	}
	return std
}

// DateListJSFormat returns chronologically sorted list of date (JS Format)
func (s *SituationToDate) DateListJSFormat() []string {
	res := make([]string, len(s.MileStones))
	var i int = 0
	for _, d := range s.MileStones {
		res[i] = d.StringJS()
		i++
	}
	sort.Strings(res)
	return res
}
