package main

import "net"

func main() {
	addr := ":2525"
	listner, err := net.Listen("tcp", addr)
	if err != nil {
	}
}
