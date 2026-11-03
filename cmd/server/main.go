package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const BUFFER_SIZE = 1024

func main() {
	conn, err := net.ListenPacket("udp", ":8080")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer conn.Close()

	b := make([]byte, BUFFER_SIZE)
	for {
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading packet: %v", err)
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
