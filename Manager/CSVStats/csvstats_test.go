package CSVStats

import (
	"fmt"
	"strings"
	"testing"
)

var csvstring string = `A;B;C
a0;a1;a2
b0;b1;b2`

func TestNewCSVStatsFrom(t *testing.T) {
	cs := NewCSVStats(NewIndexDesc("A", "A"), NewIndexDesc("AB", "A", "B"))
	err := cs.AddCSVDataFrom(strings.NewReader(csvstring))
	if err != nil {
		t.Fatal("NewCSVStatsFrom(reader) returns :", err.Error())
	}

	fmt.Printf("\nHeader :%v\nData: %v", cs, cs.data)
}
