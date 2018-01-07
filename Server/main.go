package main

import (
	"github.com/gorilla/mux"
	"github.com/lpuig/Novagile/Manager"
	"github.com/lpuig/Novagile/Route"
	"log"
	"net/http"
	"os/exec"
)

//go:generate go build -o ../server.exe

const (
	AssetsDir  = `../WebAssets`
	AssetsRoot = `/Assets/`
	RootDir    = `./Dist`

	ServicePort = ":8080"

	StatCSVFile = `./Ressources/Stats Projets Novagile.csv`
	PrjJSONFile = `./Ressources/Projets Novagile.xlsx.json`
)

func main() {
	manager, err := Manager.NewManager(PrjJSONFile, StatCSVFile)
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
	router.HandleFunc("/xls", withManager(Route.GetXLS)).Methods("GET")
	router.HandleFunc("/ptf/{prjid}", withManager(Route.UpdatePrj)).Methods("PUT")
	router.HandleFunc("/ptf/{prjid}", withManager(Route.DeletePrj)).Methods("DELETE")

	router.PathPrefix(AssetsRoot).Handler(http.StripPrefix(AssetsRoot, http.FileServer(http.Dir(AssetsDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(RootDir)))

	LaunchPageInBrowser()
	log.Print("Listening on ", ServicePort)
	log.Fatal(http.ListenAndServe(ServicePort, router))
}

func LaunchPageInBrowser() error {
	cmd := exec.Command("cmd", "/c", "start", "http://localhost:8080")
	return cmd.Start()
}

// Done Persist JSON repo after each Route request
// Done Import XLS to JSON
// Done Export JSON to XLS
// Done launch webpage with command("cmd /c start http://localhost:8080") or "explorer "http://localhost:8080""
