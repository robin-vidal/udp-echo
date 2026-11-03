package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
)

type config struct {
	port     int
	dropRate int
	bufSize  int
}

func parseFlags() *config {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "port the server listen on")
	flag.IntVar(&cfg.dropRate, "dropRate", 25, "percentage of packets to drop (0-100)")
	flag.IntVar(&cfg.bufSize, "bufferSize", 1024, "size of the message buffer")

	flag.Parse()
	return &cfg
}

func main() {
	cfg := parseFlags()
	if cfg.dropRate < 0 || cfg.dropRate > 100 {
		fmt.Fprintf(os.Stderr, "invalid drop rate %d, must be between 0 and 100\n", cfg.dropRate)
		os.Exit(1)
	}

	port := fmt.Sprintf(":%d", cfg.port)
	conn, err := net.ListenPacket("udp", port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	b := make([]byte, cfg.bufSize)
	for {
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading packet: %v\n", err)
			continue
		}

		if n > 0 {
			if rand.Intn(100) < cfg.dropRate {
				fmt.Fprintf(os.Stderr, "dropping packet from %s/n", addr)
				continue
			}

			fmt.Printf("[%s] %s\n", addr, b[:n])
			_, err := conn.WriteTo(b[:n], addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error writing packet: %v\n", err)
			}
		}
	}
}
