package logger

import (
	"log"
	"os"
)

var ErrorFileLogger *log.Logger

func InitLogger() {
	logDir := os.Mkdir("logs", 0755)
}
