package stats

import (
	"encoding/json"
	"testing"
)

func TestStat_SetValue(t *testing.T) {
	s := NewStat()
	s.SetValue("A", 0, 1)
	s.SetValue("A", 1, 2)
	s.SetValue("B", 1, 3)

	j, e := json.Marshal(s)
	if e != nil {
		t.Fatal("Stat.Marshal: " + e.Error())
	}
	println(string(j))
	s2 := NewStat()
	e = json.Unmarshal(j, &s2)
	if e != nil {
		t.Fatal("Stat.Unmarshal: " + e.Error())
	}

	if v, ok := s2.GetValue("A", 0, 0); !ok {
		t.Error("Stat.GetValue does not find serial")
		if v != 1 {
			t.Errorf("Stat.GetValue does return correct value : %f instead of 1", v)
		}
	}

	if _, ok := s2.GetValue("C", 0, 0); ok {
		t.Error("Stat.GetValue unexpectedly finds value for wrong serial")
	}

	if _, ok := s2.GetValue("B", 0, 0); ok {
		t.Error("Stat.GetValue unexpectedly finds value for wrong offset")
	}

}
