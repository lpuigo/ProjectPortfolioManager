package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
)

const (
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

func (j jirarow) QueryTo(jiraDb *sql.DB, w *csv.Writer) error {
	jiraRows, err := jiraDb.Query(jiraQuery)
	if err != nil {
		return fmt.Errorf("could not query: %s", err.Error())
	}
	defer jiraRows.Close()

	w.Write(j.Header())

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
			return fmt.Errorf("could not scan: %s", err.Error())
		}
		w.Write(j.ToStringSlice())
		//Println(j.ToStringSlice())
	}

	err = jiraRows.Err()
	if err != nil {
		log.Fatal("query returns", err)
		return fmt.Errorf("query returns: %s", err.Error())
	}
	return nil
}
