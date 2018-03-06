package main

import (
	"github.com/gorilla/mux"
	"github.com/lpuig/novagile/manager"
	"github.com/lpuig/novagile/manager/config"
	"github.com/lpuig/novagile/route"
	"log"
	"net/http"
	"os"
	"os/exec"
)

//go:generate go build -v -o ../server.exe

const (
	AssetsDir  = `../WebAssets`
	AssetsRoot = `/Assets/`
	RootDir    = `./Dist`

	ServicePort = ":8080"

	StatCSVFile = `./Ressources/Stats Projets Novagile.csv`
	PrjJSONFile = `./Ressources/Projets Novagile.xlsx.json`

	NoWebLock = true

	JiraStatDir     = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE`
	ArchivedStatDir = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Archived SRE`

	ConfigFile = `.\config.json`
	LogFile    = `.\server.log`
)

type Conf struct {
	StatInputDir   string
	StatArchiveDir string
	ServicePort    string
	LogFile        string
	NoWebLock      bool
}

func main() {

	// Init Config
	conf := &Conf{
		StatInputDir:   JiraStatDir,
		StatArchiveDir: ArchivedStatDir,
		ServicePort:    ServicePort,
		LogFile:        LogFile,
		NoWebLock:      NoWebLock,
	}
	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		log.Fatal(err)
	}

	// Init Log
	logfile := StartLog(conf.LogFile)
	defer logfile.Close()
	log.Println("Server Started =============================================================================")

	mgr, err := manager.NewManager(PrjJSONFile, StatCSVFile)
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
	router.HandleFunc("/stat/prjlist/{prjid:-?[0-9]+}", withManager(route.GetProjectStatProjectList)).Methods("GET")
	router.HandleFunc("/stat/reinit", withManager(route.GetInitProjectStat)).Methods("GET")
	router.HandleFunc("/stat/update", withManager(route.GetUpdateProjectStat)).Methods("GET")
	router.HandleFunc("/stat/{prjid:[0-9]+}", withManager(route.GetProjectStat)).Methods("GET")
	router.HandleFunc("/xls", withManager(route.GetXLS)).Methods("GET")

	router.PathPrefix(AssetsRoot).Handler(http.StripPrefix(AssetsRoot, http.FileServer(http.Dir(AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(RootDir)))

	LaunchPageInBrowser(conf.NoWebLock)
	log.Print("Listening on ", ServicePort)
	log.Fatal(http.ListenAndServe(ServicePort, router))
}

func LaunchPageInBrowser(lanchWeb bool) error {
	if lanchWeb {
		cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080")
		return cmd.Start()
	}
	log.Printf("No Web Lock found")
	return nil
}

func StartLog(logfile string) *os.File {
	//create your file with desired read/write permissions
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//set output of logs to f
	log.SetOutput(f)
	return f
}

// Done Persist JSON repo after each Route request
// Done Import XLS to JSON
// Done Export JSON to XLS
// Done launch webpage with command("cmd /c start http://localhost:8080") or "explorer "http://localhost:8080""
// Done expose import service (update stat with all csv file found in "Import" Dir, processed file are zipped and moved to "Imported" dir, or "Failed" dir if an error occurered. A file with related error is produced aside from the rejected file
// Done Create a log file containing all server activity
// TODO expose a service to upload the log file
// TODO expose an admin front end to show server activity / trigger admin operation

// TODO expose a service to show recent unmatched Jira issue
