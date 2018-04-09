package manager

import (
	"fmt"
)

type CsvFileWarning struct {
	Warnings []string
}

func NewCsvFileWarning() *CsvFileWarning {
	return &CsvFileWarning{Warnings: []string{}}
}

func (ce *CsvFileWarning) Warning() string {
	res := ""
	for _, s := range ce.Warnings {
		res += s
	}
	return res
}

func (ce *CsvFileWarning) IsError() bool {
	return len(ce.Warnings) > 0
}

func (ce *CsvFileWarning) AddWarning(numline int, signature string) *CsvFileWarning {
	ce.Warnings = append(ce.Warnings, fmt.Sprintf("line %d:%s\n", numline, signature))
	return ce
}
