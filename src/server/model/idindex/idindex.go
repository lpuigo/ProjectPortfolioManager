package idindex

type IdIndex map[int]int

func New() IdIndex {
	return IdIndex{}
}

func (i IdIndex) AddElem(id, pos int) {
	i[id] = pos
}

func (i IdIndex) DeleteElem(id int) {
	if i.IdExists(id) {
		delete(i, id)
	}
}

func (i IdIndex) IdExists(id int) bool {
	_, found := i[id]
	return found
}

func (i IdIndex) ById(id int) (int, bool) {
	e, found := i[id]
	return e, found
}
