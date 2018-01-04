package Manager

import (
	"log"
	"testing"
)

const (
	XLSPlanningFile   = `C:\Users\Laurent\Google Drive\Golang\src\github.com\lpuig\Novagile\Ressources\Projets Novagile.xlsx`
	JSONPlanningFile  = XLSPlanningFile + ".json"
	JSONPlanningFile2 = XLSPlanningFile + "2.json"
)

func TestUpdatePortfolioFromXLS(t *testing.T) {
	log.Println("Test Started")
	ptf := NewPrjPortfolio()

	err := UpdatePortfolioFromXLS(ptf, XLSPlanningFile)
	if err != nil {
		t.Error("UpdatePortfolioFromXLS returns error", err.Error())
		return
	}
	log.Println("projects loaded :", len(ptf.Projects))
	s := ptf.String()

	err = ptf.WriteJsonFile(JSONPlanningFile)
	if err != nil {
		t.Error("WriteJsonFile returns error", err.Error())
		return
	}

	var ptf2 *PrjPortfolio
	ptf2, err = NewPrjPortfolioFromJSONFile(JSONPlanningFile)
	if err != nil {
		t.Error("NewPrjPortfolioFromJSONFile returns err", err.Error())
		return
	}
	//TODO Fails because of new date compared to ref file JSONPlanningFile : find a way to avoid this
	if s != ptf2.String() {
		t.Error("UpdatePortfolioFromXLS returns improper PrjPortfolio", s)
	}
}
