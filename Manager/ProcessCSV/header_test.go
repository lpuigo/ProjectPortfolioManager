package ProcessCSV

import "testing"

func TestHeader_NewRecordSelector(t *testing.T) {
	abc := []string{"a", "b", "c", "e"}
	h := NewHeader(abc)

	sel_bc, err := h.NewRecordSelector("b", "d")
	if err == nil {
		t.Fatal("Header.NewRecordSelector does not return expected columns error :", err.Error())
	}

	sel_bc, err = h.NewRecordSelector("b", "e")
	if err != nil {
		t.Fatal("Header.NewRecordSelector returns unexpected error", err.Error())
	}

	bc := sel_bc(abc)
	if bc != "be" {
		t.Errorf("RecordSelector returns %s instead of 'bc'", bc)
	}
}
