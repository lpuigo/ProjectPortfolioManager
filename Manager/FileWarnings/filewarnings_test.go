package FileWarnings

import (
	"strings"
	"testing"
)

func TestFileWarnings_AddWarning(t *testing.T) {
	fw := NewFileWarning()

	if fw.HasWarnings() {
		t.Error("FileWarnings.HasWarnings() returns true whereas FileWarnings is empty")
	}

	fw.AddWarning(0, "warn0")
	if !fw.HasWarnings() {
		t.Error("FileWarnings.HasWarnings() returns false whereas FileWarnings is not empty")
	}

	prefix := "prefix :"
	res := fw.Warning(prefix)

	if !strings.HasPrefix(res, prefix) {
		t.Errorf("FileWarnings.Warning(%s) returns string not correctly prefixed : %s", prefix, res)
	}
}
