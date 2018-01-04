package Model

type Situations struct {
	Stds []*SituationToDate `json:"stds"`
}

func NewSituations() Situations {
	s := Situations{}
	s.Stds = make([]*SituationToDate, 0)
	return s
}

func (s Situations) String() string {
	res := "Situation : [\n"
	for _, std := range s.Stds {
		res += std.String()
	}
	res += "]\n"
	return res
}

func (s *Situations) Len() int {
	return len(s.Stds)
}

/*
func (s *Situations) Less(i, j int) bool {
	return s.Stds[i].UpdateOn.After(s.Stds[j].UpdateOn)
}

func (s *Situations) Swap(i, j int) {
	s.Stds[i], s.Stds[j] = s.Stds[j], s.Stds[i]
}
*/

// GetSituationToDate returns a new SituationToDate with the all the Latest Situation' date
func (s *Situations) GetSituationToDate() *SituationToDate {
	if s.Len() == 0 {
		return NewSituationToDate()
	}

	// get oldest situationtodate
	std := s.Stds[s.Len()-1].Clone()
	if s.Len() > 1 {
		// update it with all newer situationtodate
		for i := s.Len() - 2; i >= 0; i-- {
			std.UpdateWith(s.Stds[i])
		}
	}
	return std
}

// Push adds to current Situations the given newest SituationToDate (if nil, it is a no-op)
func (s *Situations) push(std *SituationToDate) {
	if std != nil {
		s.Stds = append([]*SituationToDate{std}, s.Stds...)
	}
}

// Update updates the current Situations s with given SituationToDate
//
// All std Milestones updated compared to s (Milestone's date changed, or added)
// will be added associated with their update date.
// All missing Milestones will be removed
func (s *Situations) Update(std *SituationToDate) {
	if std == nil {
		return
	}

	if s.Len() == 0 {
		s.Stds = append(s.Stds, std.Clone())
		return
	}

	cstd := s.GetSituationToDate()

	if std.UpdateOn.After(s.Stds[0].UpdateOn) {
		// if std is more recent than cstd
		s.push(cstd.DifferenceWith(std))
	} else if std.UpdateOn.Equal(s.Stds[0].UpdateOn) {
		// if std has same updateon date than cstd
		s.Stds[0].UpdateWith(cstd.DifferenceWith(std))
	} else {
		panic("given SituationToDate is older than current SituationToDate")
	}
}
