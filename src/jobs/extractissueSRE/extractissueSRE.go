package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lpuig/novagile/src/server/logger"
	"github.com/lpuig/novagile/src/server/manager/config"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	DateFormat = "2006-01-02"
	ConfigFile = "./extractIssueSRE.config"

	LogFile          = "./extractIssueSRE.log"
	OutputDir        = "."
	OutputFileFormat = "extract__DATE_.csv"

	UserPwd = "username:password"
	DBName  = "dbserver/basename"

	CallServer = true
	ServerURL  = "http://localhost:8080/stat/update"
)

type Conf struct {
	LogFile string

	UserPwd string
	DBName  string

	OutputDir        string
	OutputFileFormat string

	CallServer bool
	ServerUrl  string
}

//go:generate go build -v -o ../extractIssueSRE.exe

type Querier interface {
	QueryTo(db *sql.DB, w *csv.Writer) error
}

func main() {
	conf := &Conf{
		UserPwd:          UserPwd,
		DBName:           DBName,
		LogFile:          LogFile,
		OutputDir:        OutputDir,
		OutputFileFormat: OutputFileFormat,
		CallServer:       CallServer,
		ServerUrl:        ServerURL,
	}
	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		log.Fatal(err)
	}

	// Init Log
	logfile := logger.StartLog(conf.LogFile)
	defer logfile.Close()

	db := dbConnect(conf)
	defer db.Close()

	outputfilename := strings.Replace(conf.OutputFileFormat, "_DATE_", time.Now().Format(DateFormat), 1)
	outputfile := path.Join(conf.OutputDir, outputfilename)

	j := &jirarow{}

	if err := queryResultToCSVFile(j, db, outputfile); err != nil {
		log.Fatal(err)
	}

	// Call HTTP.Get to trigger resultfile loading
	if conf.CallServer {
		time.Sleep(time.Second) // timer to let the file persist to disk. Usefull ?
		if err := triggerExtractProcess(conf.ServerUrl); err != nil {
			log.Fatal("could not trigger extract processing", err.Error())
		}
	}
}

func dbConnect(conf *Conf) *sql.DB {
	jiraDb, err := sql.Open(
		"mysql",
		//"UserPwd"@"DbName",
		conf.UserPwd+"@"+conf.DBName,
	)
	if err != nil {
		log.Fatal("could not open DB", err)
	}
	return jiraDb
}

func queryResultToCSVFile(q Querier, db *sql.DB, file string) error {
	f, err := os.Create(file)
	if err != nil {
		fmt.Errorf("could not create extract file :%s", err.Error())
	}
	defer f.Close()
	log.Println("create output file:", file)

	cw := csv.NewWriter(f)
	cw.UseCRLF = true
	cw.Comma = ';'
	defer cw.Flush()

	if err := q.QueryTo(db, cw); err != nil {
		// if query fails, remove the created file
		os.Remove(file)
		return err
	}
	return nil
}

func triggerExtractProcess(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	servresp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Printf("Server replies:%s", servresp)
	return nil
}
