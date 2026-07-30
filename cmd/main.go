package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"smtp-server/internal/logger"
)

// main server initialization
func main() {
	logger.InitLogger()
	addr := ":2525"
	listner, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("error creating lister: %v", err)
		logger.ErrorFileLogger.Printf("error creating listner: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Server running on port %v", addr)
	logger.SuccessFileLogger.Printf("Server running on port: %v", addr)

	for {
		conn, err := listner.Accept()
		if err != nil {
			logger.ErrorFileLogger.Printf("Error accepting connection %v", err)
			continue
		}
		// creates thread for each connection
		go handleClient(conn)
	}
}

// handleClient([net.Conn]) handlesEvery client
func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.Write([]byte("220 whatever.com smtp Server Ready\r\n"))

	var (
		from    string
		rcpt    []string
		rawData strings.Builder
	)
	inDataMode := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			logger.ErrorFileLogger.Printf("Error reading string: %v", err)
			return
		}
		line = strings.TrimRight(line, "\r\n")
		if inDataMode {
			if line == "." {
				inDataMode = false
				conn.Write([]byte("message accepted for delivery\r\n"))
				saveAndPrintEmail(from, rcpt, rawData)
				continue
			}
			rawData.WriteString(line + "\r\n")
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		default:
			conn.Write([]byte("500 5.5.1 command unrecognized\r\n"))
			return
		}
	}
}

func saveAndPrintEmail(from string, rcpt []string, rawData any) {}
