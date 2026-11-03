package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"time"
)

type config struct {
	addr    string
	port    int
	msg     string
	timeout int
	retries int

	bufSize int
}

func parseFlags() *config {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", "localhost", "address the server listen on")
	flag.IntVar(&cfg.port, "port", 8080, "port the server listen on")
	flag.StringVar(&cfg.msg, "msg", "Hello World!", "message to send to the server")
	flag.IntVar(&cfg.timeout, "timeout", 2000, "read/write timeouts in milliseconds")
	flag.IntVar(&cfg.retries, "retries", 3, "number of retries after a timeout")
	flag.IntVar(&cfg.bufSize, "bufferSize", 1024, "size of the message buffer")

	flag.Parse()
	return &cfg
}

func main() {
	cfg := parseFlags()
	host := net.JoinHostPort(cfg.addr, strconv.Itoa(cfg.port))

	conn, err := net.Dial("udp", host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	msg := []byte(cfg.msg)
	b := make([]byte, cfg.bufSize)
	var n int

	for attempt := 0; attempt <= cfg.retries; attempt++ {
		conn.SetWriteDeadline(time.Now().Add(time.Duration(cfg.timeout) * time.Millisecond))
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "attempt %d: unable to send the message: %v\n", attempt+1, err)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.timeout) * time.Millisecond))
		n, err = conn.Read(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "attempt %d: unable to read the message: %v\n", attempt+1, err)
			continue
		}

		break
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed after %d retries\n", cfg.retries)
		os.Exit(1)
	}

	if n != len(msg) || !slices.Equal(b[:n], msg) {
		fmt.Fprintf(os.Stderr, "invalid response content:\n\tgot: %s\n\twant: %s\n", b, msg)
		os.Exit(1)
	}

	fmt.Printf("[%s] %s\n", host, b[:n])
}
