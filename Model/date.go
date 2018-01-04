package Model

import (
	"time"
)

type Date time.Time

const TimeStringLayout string = "02/01/2006"
const TimeJSLayout string = "2006-01-02"

func DateFromJSONString(s string) (Date, error) {
	date, err := time.Parse(`"`+TimeJSLayout+`"`, s)
	return Date(date), err
}

func DateFromJSString(s string) (Date, error) {
	date, err := time.Parse(TimeJSLayout, s)
	return Date(date), err
}

func DateFromString(s string) (Date, error) {
	date, err := time.Parse(TimeStringLayout, s)
	return Date(date), err
}

func (d Date) MarshalJSON() ([]byte, error) {
	str := d.ToTime().Format(`"` + TimeJSLayout + `"`)
	return []byte(str), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	//date, err := time.Parse(`"`+TimeJSLayout+`"`, string(b))
	date, err := DateFromJSONString(string(b))
	if err != nil {
		return err
	}
	*d = Date(date)
	return nil
}

func ZeroDate() Date {
	return Date(time.Time{})
}

func (d Date) ToTime() time.Time {
	return time.Time(d)
}

func (d Date) String() string {
	return d.ToTime().Format(TimeStringLayout)
}

func (d Date) StringJS() string {
	return d.ToTime().Format(TimeJSLayout)
}

func (d Date) GetMonday() Date {
	wd := int(d.ToTime().Weekday())
	if wd == 0 {
		wd = 7
	}
	wd--
	return Date(d.ToTime().AddDate(0, 0, -wd))
}

func (d Date) AddDays(n int) Date {
	return Date(d.ToTime().AddDate(0, 0, n))
}

func (d Date) DaysSince(d2 Date) int {
	return int(d.ToTime().Sub(d2.ToTime()) / time.Duration(24*time.Hour))
}

func (d Date) After(d2 Date) bool {
	return d.ToTime().After(time.Time(d2))
}

func (d Date) Before(d2 Date) bool {
	return d.ToTime().Before(time.Time(d2))
}

func (d Date) Equal(d2 Date) bool {
	return d.ToTime().Equal(time.Time(d2))
}

func (d Date) IsZero() bool {
	return d.ToTime().IsZero()
}

func Today() Date {
	return Date(time.Now().Truncate(24 * time.Hour))
}
