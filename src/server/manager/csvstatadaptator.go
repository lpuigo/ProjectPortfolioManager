package manager

import (
	"encoding/csv"
	"fmt"
	fw "github.com/lpuig/prjptf/src/server/manager/filewarnings"
	"github.com/lpuig/prjptf/src/server/model"
	"io"
	"os"
	"strconv"
)

func projectKey(client, name string) string {
	return client + "!" + name
}

func createProjectKeyCatalog(pptf *PrjPortfolio) map[string]int {
	res := map[string]int{}
	for _, p := range pptf.Projects {
		res[projectKey(p.Client, p.Name)] = p.Id
	}
	return res
}

func parseFloat(s string, numline int, csvFileWarning *fw.FileWarnings) (float64, error) {
	nf, err := strconv.ParseFloat(s, 64)
	if err != nil {
		csvFileWarning.AddWarning(numline, fmt.Sprintf("Malformed float '%s'", s))
		return 0, err
	}
	return nf, nil
}

func parseRecord(record []string, pptf *PrjPortfolio, sptf *StatPortfolio, colindex map[string]int, numline int, csvFileWarning *fw.FileWarnings) int {
	var spent, remaing, estimated float64

	prjKeyCatalog := createProjectKeyCatalog(pptf)

	pk := projectKey(record[colindex["CLIENT_NAME"]], record[colindex["PROJECT"]])
	d, err := model.DateFromString(record[colindex["DATE"]])
	if err != nil {
		csvFileWarning.AddWarning(numline, err.Error())
		return 0
	}
	d = d.GetMonday()
	id, found := prjKeyCatalog[pk]
	if !found {
		csvFileWarning.AddWarning(numline, "no project found for "+pk)
		return 0
	}

	if spent, err = parseFloat(record[colindex["SPENT"]], numline, csvFileWarning); err != nil {
		return 0
	}
	if remaing, err = parseFloat(record[colindex["REMAIN_TIME"]], numline, csvFileWarning); err != nil {
		return 0
	}
	if estimated, err = parseFloat(record[colindex["INIT_ESTIMATE"]], numline, csvFileWarning); err != nil {
		return 0
	}
	//project id found, add its stat to sptf
	sp := sptf.GetStatById(id)
	if sp == nil {
		// related project does not has stat yet, let's create it
		sp = model.NewProjectStat()
		sp.Id = id
		sp.StartDate = d
		sptf.AddProjectStat(sp)
	}
	sp.AddValues(d, spent, remaing, estimated)
	return 1
}

// UpdateStatPortfolioFromCSVFile adds correct formated stats found in given file to sptf and returns howmany stats were added, and warnings (skipped lines)
func UpdateStatPortfolioFromCSVFile(file string, pptf *PrjPortfolio, sptf *StatPortfolio) (int, error, *fw.FileWarnings) {
	csvFileWarning := fw.NewFileWarning()
	numadded := 0
	f, err := os.Open(file)
	if err != nil {
		return 0, err, nil
	}
	defer f.Close()
	csvr := csv.NewReader(f)
	csvr.Comma = ';'
	csvr.Comment = '#'

	colindex := map[string]int{}
	numline := 0
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err, nil
		}
		if numline == 0 {
			for i, v := range record {
				colindex[v] = i
			}
		} else {
			numadded += parseRecord(record, pptf, sptf, colindex, numline, csvFileWarning)
		}

		numline++
	}
	return numadded, nil, csvFileWarning

}
