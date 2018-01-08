package RecordSet

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var (
	csvstring string = `A;B;C
a0;a1;a2
b0;b1;b2`
	csvstring2 string = `A;B;C;D
a0;a1;a2;a3
b0;b1;b2;b3`
)

func testAddCSVDataFrom(cs *RecordSet, r io.Reader, t *testing.T) {
	err := cs.AddCSVDataFrom(strings.NewReader(csvstring))
	if err != nil {
		t.Fatal("AddCSVDataFrom(reader) returns :", err.Error())
	}
}

func TestNewRecord_AddCSVDataFrom(t *testing.T) {
	cs := NewRecordSet()
	testAddCSVDataFrom(cs, strings.NewReader(csvstring), t)

	expectedHeaderRecord := Record{"A", "B", "C"}
	if !cs.header.keys.Equals(expectedHeaderRecord) {
		t.Errorf("Header not properly set : %s instead of %s", cs.header.keys, expectedHeaderRecord)
	}
	if cs.Len() != 2 {
		t.Errorf("Unexpected Length : %d instead of 2", cs.Len())
	}
	expectedRecord0 := Record{"a0", "a1", "a2"}
	if !cs.Get(0).Equals(expectedRecord0) {
		t.Errorf("Unexpected record 0 : %s instead of %s", cs.Get(0), expectedRecord0)
	}
	expectedRecord1 := Record{"b0", "b1", "b2"}
	if !cs.Get(1).Equals(expectedRecord1) {
		t.Errorf("Unexpected record 1 : %s instead of %s", cs.Get(1), expectedRecord1)
	}
}

func TestRecordSet_AddCSVDataFrom2(t *testing.T) {
	cs := NewRecordSet()
	testAddCSVDataFrom(cs, strings.NewReader(csvstring), t)
	testAddCSVDataFrom(cs, strings.NewReader(csvstring), t)
	if cs.Len() != 4 {
		t.Errorf("Unexpected Length : %d instead of 4", cs.Len())
	}
}

func TestRecordSet_AddCSVDataFrom3(t *testing.T) {
	cs := NewRecordSet()
	testAddCSVDataFrom(cs, strings.NewReader(csvstring), t)
	err := cs.AddCSVDataFrom(strings.NewReader(csvstring2))
	if err == nil {
		t.Error("AddCSVDataFrom(reader) returns no error. Expected 'Uncompatible RecordSet structure : [A B C D] vs [A B C]'")
	}
}

func TestNewRecordSet_CreateSubRecordSet(t *testing.T) {
	cs := NewRecordSet()
	testAddCSVDataFrom(cs, strings.NewReader(csvstring), t)
	rs := cs.CreateSubRecordSet()
	expectedHeaderRecord := Record{"A", "B", "C"}
	if !rs.header.keys.Equals(expectedHeaderRecord) {
		t.Errorf("Header not properly set : %s instead of %s", rs.header.keys, expectedHeaderRecord)
	}
	if rs.Len() != 0 {
		t.Errorf("Unexpected Length : %d instead of 0", cs.Len())
	}
}

func TestRecordSet_WriteCSVTo(t *testing.T) {
	cs := NewRecordSet()
	testAddCSVDataFrom(cs, strings.NewReader(csvstring), t)
	f, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Fatalf("TempFile returns ", err.Error())
	}
	defer os.Remove(f.Name())
	err = cs.WriteCSVTo(f)
	if err != nil {
		t.Fatalf("WriteCSVTo returns ", err.Error())
	}

	cs2 := NewRecordSet()
	testAddCSVDataFrom(cs2, f, t)

	info := Record{"Header", "Record1", "Record2"}
	i1 := []Record{cs.header.keys, cs.Get(0), cs.Get(1)}
	i2 := []Record{cs2.header.keys, cs2.Get(0), cs2.Get(1)}

	for i, r := range i1 {
		if !r.Equals(i2[i]) {
			t.Errorf("Discrepancy on %s, %s instead of %s", info[i], r, i2[i])
		}
	}
}
