package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	fm "github.com/lpuig/novagile/src/client/frontmodel"
	mgr "github.com/lpuig/novagile/src/server/manager"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MgrHandlerFunc func(*mgr.Manager, http.ResponseWriter, *http.Request)

//var Mgr *Manager.Manager

func addError(w http.ResponseWriter, logmsg *string, errmsg string, code int) {
	*logmsg += fmt.Sprintf("%s (%d)", errmsg, code)
	http.Error(w, errmsg, code)
}

func formatLog(t time.Time, msg *string) {
	log.Printf("%s (served in %v)\n", *msg, time.Since(t))
}

func GetPtf(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request GetPtf Received from '" + r.Header.Get("Origin") + "' : "
	defer formatLog(time.Now(), &logmsg)

	w.Header().Set("Content-Type", "application/json")
	mgr.GetPrjPtf(w)
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func UpdatePrj(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request UpdatePrj Received from '" + r.Header.Get("Origin") + "' : "
	defer formatLog(time.Now(), &logmsg)

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
	defer formatLog(time.Now(), &logmsg)
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
	defer formatLog(time.Now(), &logmsg)
	defer r.Body.Close()
	vars := mux.Vars(r)
	prjid, err := strconv.Atoi(vars["prjid"])
	if err != nil {
		addError(w, &logmsg, "misformated project id '"+vars["prjid"]+"'", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = mgr.GetProjectStatProjectList(prjid, w)
	if err != nil {
		addError(w, &logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func GetInitProjectStat(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := "Request GetInitProjectStat Received from '" + r.Header.Get("Origin") + "' : "
	defer formatLog(time.Now(), &logmsg)
	defer r.Body.Close()
	w.Header().Set("Content-Type", "text/plain")
	err := mgr.ReinitStats(w)
	if err != nil {
		addError(w, &logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func GetUpdateProjectStat(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logmsg := "Request GetUpdateProjectStat Received from '" + r.Header.Get("Origin") + "' : "
	defer formatLog(time.Now(), &logmsg)
	defer r.Body.Close()
	w.Header().Set("Content-Type", "text/plain")
	err := mgr.UpdateWithNewStatFiles(w)
	if err != nil {
		addError(w, &logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func CreatePrj(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request CreatePrj Received from '" + r.Header.Get("Origin") + "' : "
	defer formatLog(time.Now(), &logmsg)

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
	defer formatLog(time.Now(), &logmsg)

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
	defer formatLog(time.Now(), &logmsg)

	w.Header().Set("Content-Disposition", "attachment; filename=\"Projet Novagile.xlsx\"")
	w.Header().Set("Content-Type", "application/vnd.ms-excel")

	mgr.GetProjectsPtfXLS(w)
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}

func GetJiraTeamLogs(mgr *mgr.Manager, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logmsg := "Request GetJiraTeamLogs Received from '" + r.Header.Get("Origin") + "' : "
	defer formatLog(time.Now(), &logmsg)

	w.Header().Set("Content-Type", "application/json")

	err := mgr.GetJiraTeamLogs(w)
	if err != nil {
		addError(w, &logmsg, err.Error(), http.StatusInternalServerError)
		return
	}
	logmsg += fmt.Sprintf("ok (%d)", http.StatusOK)
}
