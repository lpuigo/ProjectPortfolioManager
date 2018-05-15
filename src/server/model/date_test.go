package model

import (
	"testing"
	"time"
)

const (
	ZeroDateString string = "01/01/0001"
	DateString1    string = "27/11/1971"
	DateString2    string = "11/12/1972"
	DateString3    string = "24/12/1972"
	DateString4    string = "26/11/2017"
	DateString5    string = "25/11/2017"
	DateString6    string = "20/11/2017"
	ZeroDatejson   string = `"0001-01-01"`
	Datejson1      string = `"1971-11-27"`
	Datejson2      string = `"1972-12-11"`
	Datejson3      string = `"1972-12-24"`
)

func todate(s string) Date {
	date, _ := time.Parse(TimeStringLayout, s)
	return Date(date)
}

func TestZeroDate(t *testing.T) {
	zeroDate := ZeroDate()
	s := zeroDate.String()
	if s != ZeroDateString {
		t.Error("Date_ZeroDate returns", s, "instead of ", ZeroDateString)
	}

	if zeroDate.IsZero() == false {
		t.Error("Date_ZeroDate.IsZero() returns false")
	}

	zjson, err := zeroDate.MarshalJSON()
	if err != nil {
		t.Error("Date_ZeroDate.MarshalJSON() returns error", err.Error())
	}
	if string(zjson) != ZeroDatejson {
		t.Error("Date_ZeroDate.MarshalJSON() returns", string(zjson))
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	mdate, err := todate(DateString1).MarshalJSON()
	if err != nil {
		t.Error("Date_MarshalJSON returns", err)
	}
	smdate := string(mdate)
	if smdate != Datejson1 {
		t.Error("Date_MarshalJSON returns", smdate, "instead of", Datejson1)
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	d := Date{}
	err := d.UnmarshalJSON([]byte(Datejson1))
	if err != nil {
		t.Error("Date_UnmarshalJSON returns", err)
	}
	rd := todate(DateString1)
	if !d.Equal(rd) {
		t.Error("Date_UnmarshalJSON returns", time.Time(d), "instead of", time.Time(rd))
	}
}

func TestDate_After(t *testing.T) {
	if todate(DateString2).After(todate(DateString1)) != true {
		t.Error("Date_After: ", DateString2, ".After(", DateString1, ") returns false")
	}
}

func TestDate_Before(t *testing.T) {
	if todate(DateString1).Before(todate(DateString2)) != true {
		t.Error("Date_Before: ", DateString1, ".Before(", DateString2, ") returns false")
	}
}

func TestDate_Equal(t *testing.T) {
	if todate(DateString1).Equal(todate(DateString1)) != true {
		t.Error("Date_Equal: ", DateString1, ".Before(", DateString1, ") returns false")
	}
}

func TestDate_GetMonday(t *testing.T) {
	r1 := todate(DateString4).GetMonday()
	if !r1.Equal(todate(DateString6)) {
		t.Error("Date_GetMonday: ", DateString4, ".GetMonday() returns", r1.String(), "instead of", DateString6)
	}

	r2 := todate(DateString5).GetMonday()
	if !r2.Equal(todate(DateString6)) {
		t.Error("Date_GetMonday: ", DateString5, ".GetMonday() returns", r2.String(), "instead of", DateString6)
	}

	r3 := todate(DateString6).GetMonday()
	if !r3.Equal(todate(DateString6)) {
		t.Error("Date_GetMonday: ", DateString6, ".GetMonday() returns", r3.String(), "instead of", DateString6)
	}
}

func TestDate_AddDays(t *testing.T) {
	d1 := todate(DateString4)
	nd1 := d1.AddDays(-1)
	d2 := todate(DateString5)
	if !nd1.Equal(d2) {
		t.Errorf("Date_AddDays: '%s'.AddDays(-1) returns '%s' instead of '%s'", d1.String(), nd1.String(), d2.String())
	}
}

func TestDate_DaysSince(t *testing.T) {
	d1 := todate(DateString4)
	d2 := todate(DateString5)
	diff := d1.DaysSince(d2)
	if diff != 1 {
		t.Errorf("Date_DaysSince: '%s'.DaysSince('%s') returns %d instead of 1", d1.String(), d2.String(), diff)
	}
}

func TestDate_OpenDaysSince(t *testing.T) {
	for _, c := range []struct{
		beg string
		end string
		expect int
	}{
		{"14/05/2018", "21/05/2018", 5},
		{"14/05/2018", "20/05/2018", 6},
		{"14/05/2018", "23/05/2018", 7},
		{"11/05/2018", "25/05/2018", 10},
	} {
		d1 := todate(c.beg)
		d2 := todate(c.end)
		diff := d2.OpenDaysSince(d1)
		if diff != c.expect {
			t.Errorf("Date_OpenDaysSince: '%s'.OpenDaysSince('%s') returns %d instead of %d", d2.String(), d1.String(), diff, c.expect)
		}
	}
}
