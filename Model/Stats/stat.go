package Stats

type (
	timevalues map[int]float64

	Stat struct {
		Values map[string]timevalues `json:"values"`
	}
)

func newTimeValues() timevalues {
	return timevalues{}
}

func NewStat() *Stat {
	return &Stat{Values: map[string]timevalues{}}
}

func (s *Stat) SetValue(serial string, offset int, value float64) {
	tv, found := s.Values[serial]
	if !found {
		tv = newTimeValues()
		s.Values[serial] = tv
	}
	tv[offset] = value
}

func (s *Stat) GetValue(serial string, offset int, defval float64) (float64, bool) {
	tv, found := s.Values[serial]
	if !found {
		return defval, false
	}
	f, found := tv[offset]
	if !found {
		return defval, false
	}
	return f, true
}
