package FileWarnings

import (
	"fmt"
)

type FileWarnings struct {
	Warnings []string
}

func NewFileWarning() *FileWarnings {
	return &FileWarnings{Warnings: []string{}}
}

func (ce *FileWarnings) Warning(prefix string) string {
	res := ""
	for _, s := range ce.Warnings {
		res += prefix + s
	}
	return res
}

func (ce *FileWarnings) HasWarnings() bool {
	return len(ce.Warnings) > 0
}

func (ce *FileWarnings) AddWarning(numline int, signature string) *FileWarnings {
	ce.Warnings = append(ce.Warnings, fmt.Sprintf("line %d:%s\n", numline, signature))
	return ce
}
