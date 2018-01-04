package Model

import (
	"encoding/json"
	"testing"
)

const (
	SIT1Json   = `{"stds":[` + STD1json + `]}`
	SIT1String = "Situation : [\n" + STD1String + "]\n"
)

func makeSIT1() Situations {
	sit := NewSituations()
	sit.push(makeSTD1())
	return sit
}

func TestSituations_Marshal(t *testing.T) {
	sit := makeSIT1()
	b, err := json.Marshal(sit)
	if err != nil {
		t.Error("Situations Marshal returns", err)
	}
	sb := string(b)
	if sb != SIT1Json {
		t.Error("Marshal(Situations) returns improper result :", sb)
	}
}

func TestSituations_Unmarshal(t *testing.T) {
	var sit Situations
	err := json.Unmarshal([]byte(SIT1Json), &sit)
	if err != nil {
		t.Error("Unmarshal(Situations) returns err :", err)
	}
	if sit.String() != SIT1String {
		t.Error("Unmarshal(Situations) returns ", sit.String())
	}
}

func TestSituations_Update_nil(t *testing.T) {
	sit := makeSIT1()
	sit.Update(nil)
	if sit.String() != SIT1String {
		t.Error("Situations.Update(nil) returns different Situation\n", sit.String())
	}
}

func TestSituations_Update_empty(t *testing.T) {
	sit := NewSituations()
	sit.Update(makeSTD1())
	if sit.String() != SIT1String {
		t.Error("Situations.Update(STD1) returns different Situation\n", sit.String())
	}
}

func TestSituations_Update_NoChange(t *testing.T) {
	sit := makeSIT1()
	std := sit.GetSituationToDate()
	std.UpdateOn = todate(DateString2)
	sit.Update(std)
	if sit.String() != SIT1String {
		t.Error("Situations.Update(<same SituationToDate with different date>) returns different Situation\n", sit.String())
	}
}

func TestSituations_Update_WithChange(t *testing.T) {
	sit := makeSIT1()
	std2 := makeSTD2()
	sit.Update(std2)
	std2d := sit.Stds[0]

	if !std2d.UpdateOn.Equal(std2.UpdateOn) {
		t.Error("Situations.Update(STD2) returns improper UpdateOn SituationToDate :", std2d.UpdateOn, "instead of", std2.UpdateOn)
	}
	m1, m1found := std2d.MileStones[MlStKickoff]
	if m1found == true {
		t.Error("Situations.Update(STD2) latest STD contains unchanged MilesStone :", MlStKickoff, "with date", m1.String())
	}
	m2, m2found := std2d.MileStones[MlStRollOut]
	if m2found == false {
		t.Error("Situations.Update(STD2) latest STD does not contain changed MilesStone :", MlStRollOut)
	}
	if !m2.Equal(std2.MileStones[MlStRollOut]) {
		t.Error("Situations.Update(STD2) latest STD does not contain correct changed MilesStone Date:", MlStRollOut, m2.String(), "instead of", std2.MileStones[MlStRollOut].String())
	}
	m3, m3found := std2d.MileStones[MlStPilot]
	if m3found == false {
		t.Error("Situations.Update(STD2) latest STD does not contain added MilesStone :", MlStPilot)
	}
	if !m3.Equal(std2.MileStones[MlStPilot]) {
		t.Error("Situations.Update(STD2) latest STD does not contain correct added MilesStone Date:", MlStPilot, m3.String(), "instead of", std2.MileStones[MlStPilot].String())
	}
}

func TestSituations_Update_WithSecondChange(t *testing.T) {
	sit := makeSIT1()
	std2 := makeSTD2()
	std2p := makeSTD2p()
	sit.Update(std2)
	sit.Update(std2p)
	std2d := sit.Stds[0]

	if !std2d.UpdateOn.Equal(std2p.UpdateOn) {
		t.Error("Situations.Update(STD2) returns improper UpdateOn SituationToDate :", std2d.UpdateOn, "instead of", std2p.UpdateOn)
	}
	m1, m1found := std2d.MileStones[MlStKickoff]
	if m1found == true {
		t.Error("Situations.Update(STD2) latest STD contains unchanged MilesStone :", MlStKickoff, "with date", m1.String())
	}
	m2, m2found := std2d.MileStones[MlStRollOut]
	if m2found == false {
		t.Error("Situations.Update(STD2) latest STD does not contain changed MilesStone :", MlStRollOut)
	}
	if !m2.Equal(std2p.MileStones[MlStRollOut]) {
		t.Error("Situations.Update(STD2) latest STD does not contain correct updated MilesStone Date:", MlStRollOut, m2.String(), "instead of", std2p.MileStones[MlStRollOut].String())
	}
	m3, m3found := std2d.MileStones[MlStPilot]
	if m3found == false {
		t.Error("Situations.Update(STD2) latest STD does not contain added MilesStone :", MlStPilot)
	}
	if !m3.Equal(std2.MileStones[MlStPilot]) {
		t.Error("Situations.Update(STD2) latest STD does not contain correct added MilesStone Date:", MlStPilot, m3.String(), "instead of", std2.MileStones[MlStPilot].String())
	}
}

func TestSituations_GetSituationToDate_empty(t *testing.T) {
	sit := NewSituations()
	std := sit.GetSituationToDate()
	if len(std.MileStones) != 0 {
		t.Error("Situations.GetSituationToDate() returns non empty STD", std.String())
	}
}

func TestSituations_GetSituationToDate_default(t *testing.T) {
	sit := makeSIT1()
	std := sit.GetSituationToDate()
	if std.String() != STD1String {
		t.Error("Situations.GetSituationToDate() returns", std.String())
	}
}

func TestSituations_GetSituationToDate_WithChange(t *testing.T) {
	sit := makeSIT1()
	std2 := makeSTD2()
	sit.Update(std2)
	std := sit.GetSituationToDate()
	if std.String() != std2.String() {
		t.Error("Situations.GetSituationToDate() returns\n", std.String(), "instead of\n", std2.String())
	}
}

func TestSituations_GetSituationToDate_WithDeleteChange(t *testing.T) {
	sit := makeSIT1()
	std2 := makeSTD2d()
	sit.Update(std2)
	std := sit.GetSituationToDate()
	if std.String() != std2.String() {
		t.Error("Situations.GetSituationToDate() returns\n", std.String(), "instead of\n", std2.String())
	}
}
