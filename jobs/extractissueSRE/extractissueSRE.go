package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lpuig/novagile/logger"
	"github.com/lpuig/novagile/manager/config"
	"path"
	"strings"
	"time"
)

const (
	DateFormat = "2006-01-02"
	ConfigFile = "./extractIssueSRE.config"

	LogFile          = "./extractIssueSRE.log"
	OutputDir        = "."
	OutputFileFormat = "extract _DATE_.csv"

	UserPwd = "readAccess:queo7VaechileiNa"
	DBName  = "tcp(jira.acticall.com:3306)/JIRA"

	jiraQuery = `
SELECT DISTINCT 
  CONVERT(DATE_ADD(CURDATE(), INTERVAL 0 DAY), CHAR CHARACTER SET utf8mb4) AS "EXTRACT_DATE",
  CONVERT(prj.pname, CHAR CHARACTER SET utf8mb4) AS "PRODUCT",
  CONVERT(CONCAT(cflo_lotclt.customvalue, '!', cflo_lotprj.customvalue), CHAR CHARACTER SET utf8mb4) AS "CLIENT!PROJECT",
  CONVERT(cflo_act.customvalue, CHAR CHARACTER SET utf8mb4) AS "ACTIVITY",
  CONVERT(CONCAT(prj.pkey, '-', issue.issuenum), CHAR CHARACTER SET utf8mb4) AS "ISSUE",
  CONVERT(ROUND(CASE WHEN issue.timeoriginalestimate IS NULL OR issue.timeoriginalestimate <= 0 THEN 0 ELSE issue.timeoriginalestimate/3600 END, 4), CHAR CHARACTER SET utf8mb4) AS "INIT_ESTIMATE",
  CONVERT(ROUND(CASE WHEN issue.timespent IS NULL OR issue.timespent <= 0 THEN 0 ELSE issue.timespent/3600 END, 4), CHAR CHARACTER SET utf8mb4) AS "TIME_SPENT",
  CONVERT(ROUND(CASE WHEN issue.timeestimate IS NULL OR issue.timeestimate <= 0 THEN 0 ELSE issue.timeestimate/3600 END, 4), CHAR CHARACTER SET utf8mb4) AS "REMAIN_TIME",
  CONVERT(issue.summary, CHAR CHARACTER SET utf8mb4) AS "SUMMARY"
FROM
  jiraissue issue
  JOIN (project prj
    JOIN (nodeassociation noda
      JOIN (projectcategory prjc
      ) ON noda.sink_node_id = prjc.id AND prjc.id IN (10000, 10001, 10103, 10200, 10500, 10700)
    ) ON prj.id = noda.source_node_id AND noda.source_node_entity = 'Project' AND noda.sink_node_entity = 'ProjectCategory'
  ) ON issue.project = prj.id
  LEFT OUTER JOIN (customfieldvalue cflv_lot
    JOIN (customfieldoption cflo_lotprj
      JOIN customfieldoption cflo_lotclt ON cflo_lotprj.parentoptionid = cflo_lotclt.id AND cflo_lotclt.parentoptionid IS NULL
    ) ON cflv_lot.stringvalue = cflo_lotprj.id AND cflo_lotprj.parentoptionid IS NOT NULL
  ) ON issue.id = cflv_lot.issue AND cflv_lot.customfield = 12000
  LEFT OUTER JOIN (customfieldvalue cflv_act
    JOIN customfieldoption cflo_act ON cflv_act.stringvalue = cflo_act.id
  ) ON issue.id = cflv_act.issue AND cflv_act.customfield = 11901
-- WHERE
--  issue.updated >= SYSDATE() - INTERVAL 12 WEEK
ORDER BY
  prj.pname, CONCAT(prj.pkey, '-', issue.issuenum)
;

`
)

func jiraDbConnect(conf *Conf) *sql.DB {
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

type jirarow struct {
	EXTRACT_DATE  sql.NullString
	PRODUCT       sql.NullString
	CLIENTPROJECT sql.NullString
	ACTIVITY      sql.NullString
	ISSUE         sql.NullString
	INIT_ESTIMATE sql.NullString
	TIME_SPENT    sql.NullString
	REMAIN_TIME   sql.NullString
	SUMMARY       sql.NullString
}

func (j jirarow) ToStringSlice() []string {
	res := make([]string, 9)
	res[0] = j.EXTRACT_DATE.String
	res[1] = j.PRODUCT.String
	res[2] = j.CLIENTPROJECT.String
	res[3] = j.ACTIVITY.String
	res[4] = j.ISSUE.String
	res[5] = j.INIT_ESTIMATE.String
	res[6] = j.TIME_SPENT.String
	res[7] = j.REMAIN_TIME.String
	res[8] = j.SUMMARY.String

	return res
}

func (j jirarow) Header() []string {
	return []string{
		"EXTRACT_DATE",
		"PRODUCT",
		"CLIENT!PROJECT",
		"ACTIVITY",
		"ISSUE",
		"INIT_ESTIMATE",
		"TIME_SPENT",
		"REMAIN_TIME",
		"SUMMARY",
	}
}

type Conf struct {
	LogFile          string
	UserPwd          string
	DBName           string
	OutputDir        string
	OutputFileFormat string
}

//go:generate go build -v -o ../extractIssueSRE.exe

func main() {
	conf := &Conf{
		UserPwd: UserPwd,
		DBName:  DBName,
		LogFile: LogFile,
	}
	if err := config.SetFromFile(ConfigFile, conf); err != nil {
		log.Fatal(err)
	}

	// Init Log
	logfile := logger.StartLog(conf.LogFile)
	defer logfile.Close()

	outputfilename := strings.Replace(conf.OutputFileFormat, "_DATE_", time.Now().Format(DateFormat), 1)
	outputfile := path.Join(OutputDir, outputfilename)

	f, err := os.Create(outputfile)
	if err != nil {
		log.Fatal("could not create extract file :", err)
	}
	defer f.Close()
	log.Println("create output file:", outputfile)

	cw := csv.NewWriter(f)
	cw.UseCRLF = true
	cw.Comma = ';'

	jiraDb := jiraDbConnect(conf)
	defer jiraDb.Close()

	jiraRows, err := jiraDb.Query(jiraQuery)
	if err != nil {
		log.Fatal("could not query", err)
	}
	defer jiraRows.Close()

	j := jirarow{}

	cw.Write(j.Header())

	for jiraRows.Next() {
		err := jiraRows.Scan(
			&j.EXTRACT_DATE,
			&j.PRODUCT,
			&j.CLIENTPROJECT,
			&j.ACTIVITY,
			&j.ISSUE,
			&j.INIT_ESTIMATE,
			&j.TIME_SPENT,
			&j.REMAIN_TIME,
			&j.SUMMARY,
		)
		if err != nil {
			log.Fatal("could not scan", err)
		}
		cw.Write(j.ToStringSlice())
		//Println(j.ToStringSlice())
	}

	err = jiraRows.Err()
	if err != nil {
		log.Fatal("err returns", err)
	}
}
