package manager

import (
	"errors"
	"fmt"
	"github.com/lpuig/novagile/model"
	"github.com/tealeg/xlsx"
	"io"
)

const (
	PlanningSheetName = "Planning"
)

type ColIndex struct {
	ColNum  map[string]int
	ColName []string
}

func NewColIndex(r *xlsx.Row) *ColIndex {
	res := &ColIndex{
		ColNum:  make(map[string]int),
		ColName: []string{},
	}
	for i, c := range r.Cells {
		colname := string(c.Value)
		res.ColNum[colname] = i
		res.ColName = append(res.ColName, colname)
	}
	return res
}

func (c ColIndex) String() string {
	res := ""
	for i, n := range c.ColName {
		res += fmt.Sprintf("%2d : %s\n", i, n)
	}
	return res
}

func processRow(r *xlsx.Row, ptf *PrjPortfolio, index *ColIndex) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()

	getInfo := func(iname string) *xlsx.Cell {
		colnum, found := index.ColNum[iname]
		if !found {
			panic(fmt.Sprintf("Column '%s' is missing", iname))
		}
		return r.Cells[colnum]
	}

	prjId, err := getInfo("Id").Int()
	if err != nil {
		return err
	}
	prj := ptf.GetPrjById(prjId)
	if prj == nil {
		prj = model.NewProject()
		prj.Id = prjId
		ptf.AddPrj(prj)
	}

	prj.Client = getInfo("Client").Value
	prj.Name = getInfo("Projet").Value
	prj.LeadDev = getInfo("Lead Dev").Value
	prj.Status = getInfo("Statut").Value
	prj.Type = getInfo("Typo").Value
	prj.ForecastWL, err = getInfo("Charge Prév.").Float()
	if err != nil {
		prj.ForecastWL = 0
	}
	prj.CurrentWL, err = getInfo("Charge Conso.").Float()
	if err != nil {
		prj.CurrentWL = 0
	}
	prj.Comment = getInfo("Commentaire").Value

	std := model.NewSituationToDate()
	std.UpdateOn = model.Today()
	// Parse Cols related to described milstones
	// TODO Find a way to dynamically choose relevent Cols ID
	for _, i := range []int{8, 9, 10, 11, 12, 13, 14} {
		m := index.ColName[i]
		d, err := getInfo(m).Float()
		if err != nil {
			continue
		}
		date := model.Date(xlsx.TimeFromExcelTime(d, false))
		std.MileStones[m] = date
	}
	prj.Situation.Update(std)

	return nil
}

func UpdatePortfolioFromXLS(ptf *PrjPortfolio, file string) error {
	xlsfile, err := xlsx.OpenFile(file)
	if err != nil {
		return err
	}

	planningSheet, found := xlsfile.Sheet[PlanningSheetName]
	if found == false {
		return errors.New(fmt.Sprintf("UpdatePortfolioFromXLS : Unable to find sheet '%s'\n", PlanningSheetName))
	}

	//Read planning Column
	index := NewColIndex(planningSheet.Rows[0])
	//fmt.Println(index)

	for i, r := range planningSheet.Rows[1:] {
		if err := processRow(r, ptf, index); err != nil {
			return errors.New(fmt.Sprintf("Line %d : %s", i, err.Error()))
		}
	}

	return nil
}

var colTitles = []string{
	"Id",
	"Client",
	"Projet",
	"Lead Dev",
	"Statut",
	"Typo",
	"Charge Prév.",
	"Charge Conso.",
	"Commentaire",
	"Cadrage",
	"Kickoff",
	"Recette",
	"Formation",
	"Pilote",
	"Mep",
	"Mes",
}

func writeXLSHeaderRow(ptf *PrjPortfolio, sheet *xlsx.Sheet) map[string]int {
	horizTextStyle := xlsx.NewStyle()
	horizTextStyle.Alignment = xlsx.Alignment{
		Vertical: "center",
	}
	horizTextStyle.ApplyAlignment = true

	headerRow := sheet.AddRow()
	col := -1
	addCol := func(text string) *xlsx.Cell {
		col++
		cell := headerRow.AddCell()
		cell.Value = text
		cell.SetStyle(horizTextStyle)
		return cell
	}

	headerRow.SetHeight(54)
	colIndex := map[string]int{}
	for i, title := range colTitles {
		colIndex[title] = i
		addCol(title)
	}

	vertTextStyle := xlsx.NewStyle()
	vertTextStyle.Alignment = xlsx.Alignment{
		TextRotation: 90,
		Vertical:     "top",
	}
	vertTextStyle.ApplyAlignment = true

	ds := model.Today().AddDays(-30).GetMonday()
	de := model.Today().AddDays(85).GetMonday()
	for d := ds; d.Before(de); d = d.AddDays(7) {
		col++
		colIndex[d.String()] = col
		cell := headerRow.AddCell()
		cell.SetDate(d.ToTime())
		cell.NumFmt = "dd/mm"
		cell.SetStyle(vertTextStyle)
	}

	return colIndex
}

func writeXLSProjectRow(project *model.Project, rownum int, sheet *xlsx.Sheet, colIndex map[string]int) {
	std := project.Situation.GetSituationToDate()

	addMilestone := func(cell *xlsx.Cell, m string) {
		date := std.MileStones[m]
		time := date.ToTime()
		if !time.IsZero() {
			cell.SetDate(time)
			if col, found := colIndex[date.GetMonday().String()]; found {
				v := &(sheet.Cell(rownum, col).Value)
				e := ""
				if *v != "" {
					e = " "
				}
				sheet.Cell(rownum, col).Value += e + m[:1]
			}
		}
	}

	addWL := func(c *xlsx.Cell, wl float64) {
		if wl > 0 {
			c.SetFloatWithFormat(wl, "0.0")
		}
	}

	for _, title := range colTitles {
		cell := sheet.Cell(rownum, colIndex[title])
		switch title {
		case "Id":
			cell.SetInt(project.Id)
		case "Client":
			cell.SetString(project.Client)
		case "Projet":
			cell.SetString(project.Name)
		case "Lead Dev":
			cell.SetString(project.LeadDev)
		case "Statut":
			cell.SetString(project.Status)
		case "Typo":
			cell.SetString(project.Type)
		case "Commentaire":
			cell.SetString(project.Comment)
		case "Charge Prév.":
			addWL(cell, project.ForecastWL)
		case "Charge Conso.":
			addWL(cell, project.CurrentWL)
		case "Cadrage", "Kickoff", "Recette", "Formation", "Pilote", "Mep", "Mes":
			addMilestone(cell, title)
		}
	}
}

func writeXLSColsFormat(sheet *xlsx.Sheet, ColIndex map[string]int) {
	sheet.Col(ColIndex["Id"]).Width = 5
	sheet.Col(ColIndex["Id"]).Hidden = true
	sheet.Col(ColIndex["Client"]).Width = 24
	sheet.Col(ColIndex["Projet"]).Width = 24
	sheet.Col(ColIndex["Lead Dev"]).Width = 13
	sheet.Col(ColIndex["Statut"]).Width = 13
	sheet.Col(ColIndex["Typo"]).Width = 7.5
	sheet.Col(ColIndex["Charge Prév."]).Width = 7
	sheet.Col(ColIndex["Charge Conso."]).Width = 7
	sheet.Col(ColIndex["Cadrage"]).Width = 11
	sheet.Col(ColIndex["Kickoff"]).Width = 11
	sheet.Col(ColIndex["Recette"]).Width = 11
	sheet.Col(ColIndex["Formation"]).Width = 11
	sheet.Col(ColIndex["Pilote"]).Width = 11
	sheet.Col(ColIndex["Mep"]).Width = 11
	sheet.Col(ColIndex["Mes"]).Width = 11
	sheet.Col(ColIndex["Commentaire"]).Width = 35
	for col := len(colTitles); col < len(ColIndex); col++ {
		sheet.Col(col).Width = 4
	}
}

func WritePortfolioToXLS(ptf *PrjPortfolio, w io.Writer) error {
	xlsfile := xlsx.NewFile()
	xlsx.SetDefaultFont(11, "Calibri")
	planningSheet, err := xlsfile.AddSheet(PlanningSheetName)
	if err != nil {
		return err
	}

	colIndex := writeXLSHeaderRow(ptf, planningSheet)

	ordered(by("Client", true), by("Project", true)).Sort(ptf.Projects)

	for rownum, p := range ptf.Projects {
		writeXLSProjectRow(p, rownum+1, planningSheet, colIndex)
	}
	writeXLSColsFormat(planningSheet, colIndex)

	err = xlsfile.Write(w)
	if err != nil {
		return err
	}
	return nil
}
