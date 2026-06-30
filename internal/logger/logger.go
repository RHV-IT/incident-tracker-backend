package logger

import (
	"log"
	"os"
)

var ErrorFileLogger *log.Logger
var UpdateIncidentLogger *log.Logger

func InitLogger() {
	errFile, err := os.OpenFile("backgroundErrors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}
	ErrorFileLogger = log.New(errFile, "Async Error: ", log.Ldate|log.Ltime|log.Lshortfile)

	updateFile, err := os.OpenFile("updateIncidents.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open update log file: %v", err)
	}
	UpdateIncidentLogger = log.New(updateFile, "Update Incident: ", log.Ldate|log.Ltime|log.Lshortfile)
}