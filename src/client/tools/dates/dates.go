package date

import "time"

const TimeJSLayout string = "2006-01-02"

func New(s string) time.Time {
	t, err := time.Parse(TimeJSLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func NbDaysBetween(beg, end string) float64 {
	if beg == end {
		return 0
	}
	b := New(beg)
	e := New(end)
	return float64(e.Sub(b) / time.Duration(24*time.Hour))
}

func MinMax(date ...string) (min, max string) {
	min = "9999"
	max = "0000"
	for _, d := range date {
		if d == "" {
			continue
		}
		if d >= max {
			max = d
		}
		if d <= min {
			min = d
		}
	}
	return
}
