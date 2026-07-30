// Package logger for logging
package logger

import (
	"log"
	"os"
	"path/filepath"
)

var (
	ErrorFileLogger   *log.Logger
	SuccessFileLogger *log.Logger
)

// InitLogger creates the log directory creates the files and hands them to their respective loggers
func InitLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		log.Fatalf("failed to create log dir: %v", err)
	}
	errFilePath := filepath.Join(logDir, "errFile.log")
	errFile, err := os.OpenFile(errFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}
	ErrorFileLogger = log.New(errFile, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)

	successFilePath := filepath.Join(logDir, "successFile.log")
	successFile, err := os.OpenFile(successFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("failed to open error log file: %v", err)
	}
	SuccessFileLogger = log.New(successFile, "Success: ", log.Ldate|log.Ltime|log.Lshortfile)
}
