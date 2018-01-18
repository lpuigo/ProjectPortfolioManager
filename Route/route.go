package Route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	fm "github.com/lpuig/Novagile/Client/FrontModel"
	mgr "github.com/lpuig/Novagile/Manager"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	JiraStatDir = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE\`
)

type MgrHandlerFunc func(*mgr.Manager, http.ResponseWriter, *http.Request)

//var Mgr *Manager.Manager

func addError(w http.ResponseWriter, logmsg *string, errmsg string, code int) {
	*logmsg += fmt.Sprintf("%s (%d)", errmsg, code)
	http.Error(w, errmsg, code)
}

func GetPtf(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request GetPtf Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())

	w.Header().Set("Content-Type", "application/json")
	mgr.GetPrjPtf(w)
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func UpdatePrj(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request UpdatePrj Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())

	vars := mux.Vars(r)
	prjid, err := strconv.Atoi(vars["prjid"])
	if err != nil {
		addError(w, &logmsg, "misformated project id '"+vars["prjid"]+"'", http.StatusBadRequest)
		return
	}

	var prj = &fm.Project{}
	if r.Body == nil {
		addError(w, &logmsg, "request project missing", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(prj)
	if err != nil {
		addError(w, &logmsg, "unable to retrieve request project. "+err.Error(), http.StatusBadRequest)
		return
	}
	if prjid != prj.Id {
		addError(w, &logmsg, "URI Id does not match request project Id", http.StatusBadRequest)
		return
	}
	uprj := fm.CloneFEProject(prj)
	ptfPrj := mgr.GetPrjById(prjid)
	if ptfPrj == nil {
		addError(w, &logmsg, "request project Id not found", http.StatusNotFound)
		return
	}
	hasStat := mgr.UpdateProject(ptfPrj, uprj)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fm.CloneBEProject(ptfPrj, hasStat))
	logmsg += fmt.Sprintf("project Id %d updated (%d)", ptfPrj.Id, http.StatusOK)
}

func GetProjectStat(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := "Request GetProjectStat Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())
	defer r.Body.Close()
	vars := mux.Vars(r)
	prjid, err := strconv.Atoi(vars["prjid"])
	if err != nil {
		addError(w, &logmsg, "misformated project id '"+vars["prjid"]+"'", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = mgr.GetProjectStatById(prjid, w)
	if err != nil {
		addError(w, &logmsg, "unable to retreive Stats for project id '"+vars["prjid"]+"'", http.StatusBadRequest)
		return
	}
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func GetProjectStatProjectList(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := "Request GetProjectStatProjectList Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	mgr.GetProjectStatProjectList(w)
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func GetInitProjectStat(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := "Request GetInitProjectStat Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	mgr.ReinitStatsFromDir(JiraStatDir)
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func CreatePrj(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request CreatePrj Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())

	var prj = &fm.Project{}
	if r.Body == nil {
		addError(w, &logmsg, "request project missing", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(prj)
	if err != nil {
		addError(w, &logmsg, "unable to retrieve request project. "+err.Error(), http.StatusBadRequest)
		return
	}
	ptfPrj, hasStat := mgr.CreateProject(fm.CloneFEProject(prj))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fm.CloneBEProject(ptfPrj, hasStat))
	logmsg += fmt.Sprintf("New project Id %d added (%d)", ptfPrj.Id, http.StatusCreated)
}

func DeletePrj(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request DeletePrj Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())

	vars := mux.Vars(r)
	prjid, err := strconv.Atoi(vars["prjid"])
	if err != nil {
		addError(w, &logmsg, "misformated project id '"+vars["prjid"]+"'", http.StatusBadRequest)
		return
	}
	found := mgr.DeleteProject(prjid)
	//w.WriteHeader(http.StatusOK)
	logmsg += fmt.Sprintf("Project Id %d deleted (found : %v) (%d)", prjid, found, http.StatusCreated)
}

func GetXLS(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request GetXLS Received from '" + r.Header.Get("Origin") + "' : "
	defer func(t time.Time) { log.Printf("%s (served in %v)\n", logmsg, time.Since(t)) }(time.Now())

	w.Header().Set("Content-Disposition", "attachment; filename=\"Projet Novagile.xlsx\"")
	w.Header().Set("Content-Type", "application/vnd.ms-excel")

	mgr.GetProjectsPtfXLS(w)
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}
