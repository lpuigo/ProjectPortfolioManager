package logger

import (
	"log"
	"os"
)

// StartLog setups logger to write log in given file
//
// Do defer file.Close() just after StartLog call to ensure proper log file closing
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
