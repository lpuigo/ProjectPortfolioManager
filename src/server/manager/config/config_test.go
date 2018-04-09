package config

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	testConfigDir  = `C:\Users\Laurent\Golang\src\github.com\lpuig\Novagile\Ressources\Test`
	testConfigFile = `testConfig.json`
)

var testConf = struct {
	TestValue string
}{
	TestValue: "default value",
}

func TestSetFromFile(t *testing.T) {
	file := filepath.Join(testConfigDir, testConfigFile)
	os.Remove(file)
	err := SetFromFile(file, &testConf)
	if err != nil {
		t.Error("SetFromFile (init) returns", err.Error())
	}

	v := testConf.TestValue
	testConf.TestValue = "new value"

	err = SetFromFile(file, &testConf)
	if err != nil {
		t.Error("SetFromFile (reload) returns", err.Error())
	}
	if testConf.TestValue != v {
		t.Error("testConf is not correctly restored : %s instead of %s", testConf.TestValue, v)
	}

}
