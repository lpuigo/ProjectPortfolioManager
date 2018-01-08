package RecordSet

import "testing"

func TestHeader_NewRecordSelector(t *testing.T) {
	abc := []string{"a", "b", "c", "e"}
	h := NewHeader(abc)

	sel_bc, err := h.NewKeyGenerator("b", "d")
	if err == nil {
		t.Fatal("Header.NewKeyGenerator does not return expected columns error :", err.Error())
	}

	sel_bc, err = h.NewKeyGenerator("b", "e")
	if err != nil {
		t.Fatal("Header.NewKeyGenerator returns unexpected error", err.Error())
	}

	bc := sel_bc(abc)
	if bc != "!b!e" {
		t.Errorf("KeyGenerator returns %s instead of 'bc'", bc)
	}
}
