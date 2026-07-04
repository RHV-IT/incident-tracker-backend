package logger

import (
	"log"
	"os"
	"path/filepath"
)

var (
	ErrorFileLogger      *log.Logger
	UpdateIncidentLogger *log.Logger
	CommentLogger        *log.Logger
)

func InitLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	errFilePath := filepath.Join(logDir, "backgroundErrors.log")
	errFile, err := os.OpenFile(errFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}
	ErrorFileLogger = log.New(errFile, "Async Error: ", log.Ldate|log.Ltime|log.Lshortfile)

	updateFilePath := filepath.Join(logDir, "updateIncidents.log")
	updateFile, err := os.OpenFile(updateFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Failed to open update log file: %v", err)
	}
	UpdateIncidentLogger = log.New(updateFile, "Update Incident: ", log.Ldate|log.Ltime|log.Lshortfile)

	commentFilePath := filepath.Join(logDir, "comments.log")
	commentFile, err := os.OpenFile(commentFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Failed to open comments log file: %v", err)
	}
	CommentLogger = log.New(commentFile, "Comment: ", log.Ldate|log.Ltime|log.Lshortfile)
}
