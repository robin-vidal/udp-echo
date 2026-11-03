package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
)

const (
	BUFFER_SIZE = 1024
)

func parseFlags() (string, int, string) {
	var addr string
	var port int
	var msg string

	flag.StringVar(&addr, "addr", "localhost", "address the server listen on")
	flag.IntVar(&port, "port", 8080, "port the server listen on")
	flag.StringVar(&msg, "msg", "Hello World!", "message to send to the server")

	flag.Parse()
	return addr, port, msg
}

func main() {
	addr, port, msg := parseFlags()
	host := net.JoinHostPort(addr, strconv.Itoa(port))

	conn, err := net.Dial("udp", host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	msgBytes := []byte(msg)
	_, err = conn.Write(msgBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to send the message: %v\n", err)
		os.Exit(1)
	}

	b := make([]byte, BUFFER_SIZE)
	n, err := conn.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read the message: %v\n", err)
		os.Exit(1)
	}

	if n != len(msgBytes) || !slices.Equal(b[:n], msgBytes) {
		fmt.Fprintf(os.Stderr, "invalid response content:\n\tgot: %s\n\twant: %s\n", b, msgBytes)
		os.Exit(1)
	}

	fmt.Printf("[%s] %s\n", host, b[:n])
}
