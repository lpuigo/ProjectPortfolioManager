package RecordSet

type Record []string

func (r Record) Equals(r2 Record) bool {
	if len(r) != len(r2) {
		return false
	}
	for i, v := range r {
		if r2[i] != v {
			return false
		}
	}
	return true
}
