package main

import "testing"

func Test_triggerExtractProcess(t *testing.T) {
	err := triggerExtractProcess("http://localhost:8080/stat/update")
	if err != nil {
		t.Error(err)
	}
}
