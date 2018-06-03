package date

import "strconv"

func FormatHour(v float64) string {
	d := int(v)
	hf := (v-float64(d))*8
	h := int(hf)
	m := int((hf-float64(h))*60)
	res := ""
	if d>0 {
		res += strconv.Itoa(d)+"d"
	}
	res += strconv.Itoa(h+100)[1:]+"h"
	res += strconv.Itoa(m+100)[1:]
	return res
}