package CSVStats

import (
	"fmt"
	"strings"
	"testing"
)

var csvstring string = `A;B;C
a0;a1;a2
a0;a11;a21
b0;b1;b2`

func TestNewCSVStatsFrom(t *testing.T) {
	cs := NewCSVStats(NewIndexDesc("A", "A"), NewIndexDesc("AB", "A", "B"))
	err := cs.AddCSVDataFrom(strings.NewReader(csvstring))
	if err != nil {
		t.Fatal("NewCSVStatsFrom(reader) returns :", err.Error())
	}

	//TODO make proper tests
	fmt.Printf("\nHeader :%v\nData: %v", cs, cs.data)

	fmt.Printf("\nKeys :%v\n", cs.GetIndexKeys("A"))
	fmt.Printf("\nKeys :%v\n", cs.GetIndexKeys("AB"))
	fmt.Printf("\nRecord :%v\n", cs.GetRecordsByIndexKey("A", "!a0"))
	fmt.Printf("\nRecord :%v\n", cs.GetRecordsByIndexKey("AB", "!a0!a11"))
}

func TestCSVStats_Max(t *testing.T) {
	cs := NewCSVStats(NewIndexDesc("A", "A"), NewIndexDesc("AB", "A", "B"))
	err := cs.AddCSVDataFrom(strings.NewReader(csvstring))
	if err != nil {
		t.Fatal("NewCSVStatsFrom(reader) returns :", err.Error())
	}
	max := cs.Max("A", "!a0", NewIndexDesc("", "B"))
	t.Errorf("CSVStat.Max %v", max)
}
