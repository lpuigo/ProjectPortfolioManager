package main

import fm "github.com/lpuig/novagile/client/frontmodel"

func createStatuts() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("1 - Candidate", "Candidate"))
	res = append(res, fm.NewValText("2 - Outlining", "Outline in progress"))
	res = append(res, fm.NewValText("3 - On Going", "On Going"))
	res = append(res, fm.NewValText("4 - UAT", "UAT in progress"))
	res = append(res, fm.NewValText("5 - Monitoring", "Monitoring in progress"))
	res = append(res, fm.NewValText("6 - Done", "Done"))
	res = append(res, fm.NewValText("0 - Lost", "Candidate lost"))
	return res
}

func createTypes() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("Legacy", "Legacy"))
	res = append(res, fm.NewValText("Acti", "Novagile for Acticall"))
	res = append(res, fm.NewValText("Nov", "Novagile for Client"))
	res = append(res, fm.NewValText("Sitel", "Novagile for Sitel"))
	res = append(res, fm.NewValText("Specif", "Specific"))
	res = append(res, fm.NewValText("Archi", "Architecture"))
	res = append(res, fm.NewValText("Roadmap", "Roadmap"))
	return res
}

func createRisks() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("0", "No Risk"))
	res = append(res, fm.NewValText("1", "Low Risk"))
	res = append(res, fm.NewValText("2", "High Risk"))
	return res
}

func createMilestoneKeys() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("Kickoff", "K"))
	res = append(res, fm.NewValText("Outline", "C"))
	res = append(res, fm.NewValText("UAT", "R"))
	res = append(res, fm.NewValText("Training", "F"))
	res = append(res, fm.NewValText("Pilot End", "P"))
	res = append(res, fm.NewValText("RollOut", "M"))
	res = append(res, fm.NewValText("GoLive", "S"))
	return res
}
