package logger

import (
	"log"
	"os"
)

var ErrorFileLogger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("backgroundErrors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	ErrorFileLogger = log.New(file, "Async Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}