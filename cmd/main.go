package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/mail"
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
		to      []string
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
				saveAndPrintEmail(from, to, rawData.String())
				continue
			}
			rawData.WriteString(line + "\r\n")
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "HELO", "EHLO":
			conn.Write([]byte("200 message ready to accept message\r\n"))

		case "MAIL":
			from = strings.Trim(parts[1][5:], "<> ")
			conn.Write([]byte("250 2.1.0 Sender OK\r\n"))

		case "RCPT":
			rcpt := strings.Trim(parts[1][3:], "<> ")
			to = append(to, rcpt)
			conn.Write([]byte("250 2.1.5 Recipient OK\r\n"))

		case "DATA":
			inDataMode = true
			rawData.Reset()
			conn.Write([]byte("354 Start mail input; end with <CR><LF>,<CR><LF>\r\n"))

		case "QUIT":
			return

		default:
			conn.Write([]byte("500 5.5.1 command unrecognized\r\n"))
			return
		}
	}
}

func saveAndPrintEmail(from string, to []string, rawData string) {
	msg, err := mail.ReadMessage(strings.NewReader(rawData))
	if err != nil {
		fmt.Printf("failed to parse email: %v\n", err)
		logger.ErrorFileLogger.Printf("failed to parse email: %v\n", err)
		return
	}
	bodyBytes, _ := io.ReadAll(msg.Body)
	fmt.Println("\n==========================================")
	fmt.Printf("📬 INBOUND MAIL DELIVERED\n")
	fmt.Printf("Envelope From : %s\n", from)
	fmt.Printf("Envelope To   : %v\n", to)
	fmt.Printf("Subject       : %s\n", msg.Header.Get("Subject"))
	fmt.Printf("Date          : %s\n", msg.Header.Get("Date"))
	fmt.Println("----------------- Body -------------------")
	fmt.Println(string(bodyBytes))
	fmt.Println("==========================================")
	logger.EmailLogger.Printf(`
			\n==========================================
			📬 INBOUND MAIL DELIVERED\n
			Envelope From : %s\n
			Envelope To   : %v\n
			Subject       : %s\n
			Date          : %s\n
			----------------- Body -------------------
			%s
			==========================================
		`, from, to, msg.Header.Get("Subject"), msg.Header.Get("Date"), string(bodyBytes))
}
