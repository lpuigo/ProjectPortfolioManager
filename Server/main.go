package main

import (
	"github.com/gorilla/mux"
	"github.com/lpuig/Novagile/Manager"
	"github.com/lpuig/Novagile/Manager/Config"
	"github.com/lpuig/Novagile/Route"
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

	NoWebLockFile = `./Ressources/NoWebOpening.lock`

	JiraStatDir     = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Extract SRE`
	ArchivedStatDir = `C:\Users\Laurent\Google Drive\Travail\NOVAGILE\Gouvernance\Stat Jira\Archived SRE`

	ConfigFile = `.\config.json`
)

type Conf struct {
	StatInputDir   string
	StatArchiveDir string
}

func main() {

	config := &Conf{
		StatInputDir:   JiraStatDir,
		StatArchiveDir: ArchivedStatDir,
	}
	Config.SetFromFile(ConfigFile, config)

	manager, err := Manager.NewManager(PrjJSONFile, StatCSVFile)
	if err != nil {
		log.Fatal(err)
	}
	err = manager.AddStatFileDirs(config.StatInputDir, config.StatArchiveDir)
	if err != nil {
		log.Fatal(err)
	}

	withManager := func(hf Route.MgrHandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			hf(manager, w, r)
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/ptf", withManager(Route.GetPtf)).Methods("GET")
	router.HandleFunc("/ptf", withManager(Route.CreatePrj)).Methods("POST")
	router.HandleFunc("/ptf/{prjid:[0-9]+}", withManager(Route.UpdatePrj)).Methods("PUT")
	router.HandleFunc("/ptf/{prjid:[0-9]+}", withManager(Route.DeletePrj)).Methods("DELETE")
	router.HandleFunc("/stat/prjlist/{prjid:-?[0-9]+}", withManager(Route.GetProjectStatProjectList)).Methods("GET")
	router.HandleFunc("/stat/reinit", withManager(Route.GetInitProjectStat)).Methods("GET")
	router.HandleFunc("/stat/update", withManager(Route.GetUpdateProjectStat)).Methods("GET")
	router.HandleFunc("/stat/{prjid:[0-9]+}", withManager(Route.GetProjectStat)).Methods("GET")
	router.HandleFunc("/xls", withManager(Route.GetXLS)).Methods("GET")

	router.PathPrefix(AssetsRoot).Handler(http.StripPrefix(AssetsRoot, http.FileServer(http.Dir(AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(RootDir)))

	LaunchPageInBrowser(NoWebLockFile)
	log.Print("Listening on ", ServicePort)
	log.Fatal(http.ListenAndServe(ServicePort, router))
}

func LaunchPageInBrowser(lockfile string) error {
	_, err := os.Stat(lockfile)
	if err != nil && os.IsNotExist(err) {
		cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080")
		return cmd.Start()
	}
	log.Printf("No Web Opening Lockfile found")
	return nil
}

// Done Persist JSON repo after each Route request
// Done Import XLS to JSON
// Done Export JSON to XLS
// Done launch webpage with command("cmd /c start http://localhost:8080") or "explorer "http://localhost:8080""
// Done expose import service (update stat with all csv file found in "Import" Dir, processed file are zipped and moved to "Imported" dir, or "Failed" dir if an error occurered. A file with related error is produced aside from the rejected file
// TODO Create a log file containing all server activity
// TODO expose a service to upload the log file
// TODO expose an admin front end to show server activity / trigger admin operation

// TODO expose a service to show recent unmatched Jira issue
