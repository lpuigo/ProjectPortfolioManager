package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lpuig/prjptf/src/server/logger"
	"github.com/lpuig/prjptf/src/server/manager"
	"github.com/lpuig/prjptf/src/server/manager/config"
	"github.com/lpuig/prjptf/src/server/route"
	"log"
	"net/http"
	"os/exec"
)

//go:generate go build -v -o ../../server.exe

const (
	AssetsDir  = `../WebAssets`
	AssetsRoot = `/Assets/`
	RootDir    = `./Dist`

	ServicePort = ":8080"

	StatCSVFile = `./Ressources/Stats Projets Novagile.csv`
	PrjJSONFile = `./Ressources/Projets Novagile.xlsx.json`

	LaunchWebBrowser = true

	JiraStatDir     = `C:/Users/Laurent/Google Drive/Travail/NOVAGILE/Gouvernance/Stat Jira/Extract SRE`
	ArchivedStatDir = `C:/Users/Laurent/Google Drive/Travail/NOVAGILE/Gouvernance/Stat Jira/Archived SRE`

	ConfigFile = `./config.json`
	LogFile    = `./server.log`

	JiraUsrPwd = `usr:pwd`
	JiraDBName = `server/dbname`
)

type Conf struct {
	StatInputDir     string
	StatArchiveDir   string
	ServicePort      string
	LogFile          string
	LaunchWebBrowser bool
	JiraUsrPwd       string
	JiraDBName       string
}

func main() {

	// Init Config
	conf := &Conf{
		StatInputDir:     JiraStatDir,
		StatArchiveDir:   ArchivedStatDir,
		ServicePort:      ServicePort,
		LogFile:          LogFile,
		LaunchWebBrowser: LaunchWebBrowser,
		JiraUsrPwd:       JiraUsrPwd,
		JiraDBName:       JiraDBName,
	}
	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		log.Fatal(err)
	}

	// Init Log
	logfile := logger.StartLog(conf.LogFile)
	defer logfile.Close()
	log.Println("Server Started =============================================================================")

	mgr, err := manager.NewManager(PrjJSONFile, StatCSVFile, conf.JiraUsrPwd, conf.JiraDBName)
	if err != nil {
		log.Fatal(err)
	}
	err = mgr.AddStatFileDirs(conf.StatInputDir, conf.StatArchiveDir)
	if err != nil {
		log.Fatal(err)
	}

	withManager := func(hf route.MgrHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			hf(mgr, w, r)
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/ptf", withManager(route.GetPtf)).Methods("GET")
	router.HandleFunc("/ptf", withManager(route.CreatePrj)).Methods("POST")
	router.HandleFunc("/ptf/{prjid:[0-9]+}", withManager(route.UpdatePrj)).Methods("PUT")
	router.HandleFunc("/ptf/{prjid:[0-9]+}", withManager(route.DeletePrj)).Methods("DELETE")
	router.HandleFunc("/ptf/workload", withManager(route.GetWorkload)).Methods("GET")
	router.HandleFunc("/stat/prjlist/{prjid:-?[0-9]+}", withManager(route.GetProjectStatProjectList)).Methods("GET")
	router.HandleFunc("/stat/reinit", withManager(route.GetInitProjectStat)).Methods("GET")
	router.HandleFunc("/stat/update", withManager(route.GetUpdateProjectStat)).Methods("GET")
	router.HandleFunc("/stat/{prjid:[0-9]+}", withManager(route.GetProjectStat)).Methods("GET")
	router.HandleFunc("/jira/teamlogs", withManager(route.GetJiraTeamLogs)).Methods("GET")
	router.HandleFunc("/jira/projectlogs", withManager(route.GetJiraProjectLogs)).Methods("GET")
	router.HandleFunc("/xls", withManager(route.GetXLS)).Methods("GET")

	router.PathPrefix(AssetsRoot).Handler(http.StripPrefix(AssetsRoot, http.FileServer(http.Dir(AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(RootDir)))

	// TODO consider using github.com/nytimes/gziphandler (more versatile / efficient ?)
	gzipedrouter := handlers.CompressHandler(router)
	//gzipedrouter := router

	LaunchPageInBrowser(conf)
	log.Print("Listening on ", ServicePort)
	log.Fatal(http.ListenAndServe(ServicePort, gzipedrouter))
}

func LaunchPageInBrowser(c *Conf) error {
	if !c.LaunchWebBrowser {
		return nil
	}
	cmd := exec.Command("cmd", "/c", "start", "http://localhost"+c.ServicePort)
	return cmd.Start()
}

// Done Persist JSON repo after each Route request
// Done Import XLS to JSON
// Done Export JSON to XLS
// Done launch webpage with command("cmd /c start http://localhost:8080") or "explorer "http://localhost:8080""
// Done expose import service (update stat with all csv file found in "Import" Dir, processed file are zipped and moved to "Imported" dir, or "Failed" dir if an error occurered. A file with related error is produced aside from the rejected file
// Done Create a log file containing all server activity
// TODO Create rules in RecordSet to format record (eg : ensure SRE are formated as %.4f)
// TODO expose a service to upload the log file
// TODO expose an admin front end to show server activity / trigger admin operation

// TODO expose a service to show recent unmatched Jira issue
