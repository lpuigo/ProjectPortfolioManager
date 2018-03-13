package main

import fm "github.com/lpuig/novagile/client/frontmodel"

func createStatuts() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("1 - Candidat", "Candidat"))
	res = append(res, fm.NewValText("2 - Cadrage", "Cadrage en cours"))
	res = append(res, fm.NewValText("3 - En Cours", "Développement en cours"))
	res = append(res, fm.NewValText("4 - Recette", "Recette client en cours"))
	res = append(res, fm.NewValText("5 - SAV", "Suivi post Mise en Service"))
	res = append(res, fm.NewValText("6 - Terminé", "Terminé"))
	res = append(res, fm.NewValText("0 - Candidat Perdu", "Candidat Perdu"))
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
	res = append(res, fm.NewValText("Cadrage", "C"))
	res = append(res, fm.NewValText("Recette", "R"))
	res = append(res, fm.NewValText("Formation", "F"))
	res = append(res, fm.NewValText("Pilote", "P"))
	res = append(res, fm.NewValText("Mep", "M"))
	res = append(res, fm.NewValText("Mes", "S"))
	return res
}
