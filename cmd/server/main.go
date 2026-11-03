package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	PORT        = 8080
	BUFFER_SIZE = 1024
)

func main() {
	port := fmt.Sprintf(":%s", strconv.Itoa(PORT))
	conn, err := net.ListenPacket("udp", port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	b := make([]byte, BUFFER_SIZE)
	for {
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading packet: %v\n", err)
			continue
		}

		if n > 0 {
			fmt.Printf("[%s] %s\n", addr, b[:n])
			_, err := conn.WriteTo(b[:n], addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error writing packet: %v\n", err)
			}
		}
	}
}
