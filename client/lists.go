package main

import fm "github.com/lpuig/novagile/client/frontmodel"

func createStatuts() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("1 - Candidat", "Candidat"))
	res = append(res, fm.NewValText("2 - Cadrage", "Cadrage en cours"))
	res = append(res, fm.NewValText("3 - En Cours", "Développement en cours"))
	res = append(res, fm.NewValText("4 - Pilote", "Phase Pilote client en cours"))
	res = append(res, fm.NewValText("5 - SAV", "Suivi post Mise en Service"))
	res = append(res, fm.NewValText("6 - Terminé", "Terminé"))
	res = append(res, fm.NewValText("0 - Candidat Perdu", "Candidat Perdu"))
	return res
}

func createTypes() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("Legacy", "Legacy"))
	res = append(res, fm.NewValText("Acti", "Novagile pour Acticall"))
	res = append(res, fm.NewValText("Nov", "Novagile pour Client"))
	res = append(res, fm.NewValText("Sitel", "Novagile pour Sitel"))
	res = append(res, fm.NewValText("Specif", "Spécifique"))
	res = append(res, fm.NewValText("Archi", "Architecture"))
	return res
}

func createRisks() []*fm.ValText {
	res := []*fm.ValText{}
	res = append(res, fm.NewValText("0", "Normal"))
	res = append(res, fm.NewValText("1", "Moyen"))
	res = append(res, fm.NewValText("2", "Elevé"))
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
