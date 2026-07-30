package main

import (
	"fmt"
	"net"
	"os"

	"smtp-server/internal/logger"
)

func main() {
	addr := ":2525"
	listner, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("error creating lister: %v", err)
		logger.ErrorFileLogger.Printf("error creating listner: %v", err)
		os.Exit(1)
	}
}
