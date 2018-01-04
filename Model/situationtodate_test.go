package Model

import (
	"encoding/json"
	"testing"
)

const (
	STD1json   = `{"update":"1971-11-27","milestones":{"Kickoff":"1971-11-27","Mep":"1972-12-11"}}`
	STD1String = "SituationToDate {\n\tUpdateOn : 27/11/1971\n\tMileStones : {\n\t\tKickoff : 27/11/1971\n\t\tMep : 11/12/1972}\n}\n"
)

func makeSTD1() *SituationToDate {
	std := NewSituationToDate()
	std.UpdateOn = todate(DateString1)
	std.MileStones[MlStKickoff] = todate(DateString1)
	std.MileStones[MlStRollOut] = todate(DateString2)
	return std
}

func makeSTD2() *SituationToDate {
	std := makeSTD1()
	std.UpdateOn = todate(DateString2)
	std.MileStones[MlStPilot] = todate(DateString2)
	std.MileStones[MlStRollOut] = todate(DateString3)
	return std
}

func makeSTD2d() *SituationToDate {
	std := makeSTD2()
	delete(std.MileStones, MlStRollOut)
	return std
}

func makeSTD2p() *SituationToDate {
	std := makeSTD2()
	std.MileStones[MlStRollOut] = todate(DateString2)
	return std
}

func TestSituationToDate_Marshal(t *testing.T) {
	std := makeSTD1()
	b, err := json.Marshal(std)
	if err != nil {
		t.Error("SituationToDate Marshal returns", err)
	}
	sb := string(b)
	if sb != STD1json {
		t.Error("Marshal(SituationToDate) returns improper result :", string(sb))
	}
}

func TestSituationToDate_Unmarshal(t *testing.T) {
	var std SituationToDate
	err := json.Unmarshal([]byte(STD1json), &std)
	if err != nil {
		t.Error("Unmarshal(SituationToDate) returns err :", err)
	}
	if std.String() != STD1String {
		t.Error("Unmarshal(SituationToDate) returns ", std)
	}
}

func TestSituationToDate_String(t *testing.T) {
	std := makeSTD1()
	s := std.String()
	if s != STD1String {
		t.Error("SituationToDate.String() returns improper result :\n", s)
	}
}

func TestSituationToDate_Clone(t *testing.T) {
	std1 := makeSTD1()
	std2 := std1.Clone()

	if std1 == std2 {
		t.Error("SituationToDate.Clone() returns the same object address")
	}
	s1, s2 := std1.String(), std2.String()
	if s1 != s2 {
		t.Error("SituationToDate.Clone() returns improper cloned object :\n", s2, "\ninstead of :\n", s1)
	}
	std1.UpdateOn = todate(DateString2)
	std1.MileStones[MlStKickoff] = todate(DateString2)
	delete(std1.MileStones, MlStRollOut)
	std1.MileStones[MlStGoLive] = todate(DateString1)
	s2 = std2.String()
	if s2 != STD1String {
		t.Error("SituationToDate changed since its Clone has been updated :\n", s1, "\ninstead of :\n", STD1String)
	}
}

func TestSituationToDate_UpdateWith_null(t *testing.T) {
	std1 := makeSTD1()
	s1 := std1.String()

	std1.UpdateWith(nil)
	if std1.String() != s1 {
		t.Error("SituationToDate.UpdateWith(nil) alters SituationToDate")
	}

	std1.UpdateWith(std1.Clone())
	if std1.String() != s1 {
		t.Error("SituationToDate.UpdateWith(unchanged) alters SituationToDate")
	}
}

func TestSituationToDate_UpdateWith(t *testing.T) {
	std1, std2 := makeSTD1(), makeSTD2()
	std3 := std1.Clone()
	std3.UpdateWith(std2)
	if std3.UpdateOn.Equal(std2.UpdateOn) != true {
		t.Error("SituationToDateOld.UpdateWith(stdNew) returns UpdateOn Date differents than stdNew's")
	}

	m1, m1found := std3.MileStones[MlStKickoff]
	if m1found != true {
		t.Error("SituationToDateOld.UpdateWith(stdNew) does not return unchanged MileStone :", MlStKickoff)
	}
	if m1.Equal(std1.MileStones[MlStKickoff]) != true {
		t.Error("SituationToDateOld.UpdateWith(stdNew) alters unchanged MileStone :", MlStKickoff, "with", m1.String(), "instead of", std1.MileStones[MlStKickoff])
	}

	m2, m2found := std3.MileStones[MlStRollOut]
	if m2found != true {
		t.Error("SituationToDateOld.UpdateWith(stdNew) does not return changed MileStone :", MlStRollOut)
	}
	if !m2.Equal(std2.MileStones[MlStRollOut]) {
		t.Error("SituationToDateOld.UpdateWith(stdNew) does not return proper changed MileStone Date:", m2, "instead of", std2.MileStones[MlStRollOut])
	}

	m3, m3found := std3.MileStones[MlStPilot]
	if m3found != true {
		t.Error("SituationToDateOld.UpdateWith(stdNew) does not return added MileStone :", MlStPilot)
	}
	if !m3.Equal(std2.MileStones[MlStPilot]) {
		t.Error("SituationToDateOld.UpdateWith(stdNew) does not return proper added MileStone Date:", m3, "instead of", std2.MileStones[MlStPilot])
	}
}

func TestSituationToDate_DifferenceWith_null(t *testing.T) {
	std1 := makeSTD1()
	std2 := std1.Clone()
	std2.UpdateOn = todate(DateString2)

	sd10 := std1.DifferenceWith(nil)
	if sd10 != nil {
		t.Error("SituationToDate.DifferenceWith(nil) return non nil value", sd10.String())
	}
	sd12 := std1.DifferenceWith(std2)
	if sd12 != nil {
		t.Error("SituationToDate.DifferenceWith(<only updateon changed>) return non nil value", sd12.String())
	}
}

func TestSituationToDate_DifferenceWith(t *testing.T) {
	std1, std2 := makeSTD1(), makeSTD2()
	sd12 := std1.DifferenceWith(std2)
	if sd12.UpdateOn.Equal(std2.UpdateOn) == false {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) returns UpdateOn Date differents than stdNEw's")
	}
	m1, m1found := sd12.MileStones[MlStKickoff]
	if m1found == true {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) returns unchanged MileStone :", MlStKickoff, "=", m1.String())
	}
	m2, m2found := sd12.MileStones[MlStRollOut]
	if m2found != true {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) does not return changed MileStone :", MlStRollOut)
	}
	if !m2.Equal(std2.MileStones[MlStRollOut]) {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) does not return proper changed MileStone Date:", m2, "instead of", std2.MileStones[MlStRollOut])
	}
	m3, m3found := sd12.MileStones[MlStPilot]
	if m3found != true {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) does not return added MileStone :", MlStPilot)
	}
	if !m3.Equal(std2.MileStones[MlStPilot]) {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) does not return proper added MileStone Date:", m3, "instead of", std2.MileStones[MlStPilot])
	}
}

func TestSituationToDate_DifferenceWith_DeletedMilestone(t *testing.T) {
	std1, std2 := makeSTD1(), makeSTD2d()
	sd12 := std1.DifferenceWith(std2)

	d, found := sd12.MileStones[MlStRollOut]
	if !found {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) does not return deleted MileStone :", MlStRollOut)
		return
	}
	if !d.IsZero() {
		t.Error("SituationToDateOld.DifferenceWith(stdNew) returns deleted MileStone with not ZeroDate:", d.String())
	}
}

func TestSituationToDate_UpdateWith_DifferenceWith_DeletedMilestone(t *testing.T) {
	std1, std2 := makeSTD1(), makeSTD2d()
	std1.UpdateWith(std1.DifferenceWith(std2))
	diff := std1.DifferenceWith(std2)

	if diff != nil {
		t.Error("std1.UpdateWith(std1.DifferenceWith(std2)) does not make std1 same as std2")
	}
}
