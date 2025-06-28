package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	address := "localhost:42069"
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("ERROR: Could not resolve %s: %s", address, err)
	}
	connection, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("ERROR: Could not establish UDP connection: %s", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	close := []byte{}
	close = fmt.Append(close, "### Connection closed\n")
	go func() {
		<-signals
		fmt.Println("\nShutting down...")
		_, err = connection.Write(close)
		if err != nil {
			fmt.Printf("ERROR: Couldn't write data to connection: %s", err)
		}
		connection.Close()
		os.Exit(0)
	}()

	out := []byte{}
	out = fmt.Appendf(out, "### Connection from %d:%d\n", addr.IP, addr.Port)
	_, err = connection.Write(out)
	if err != nil {
		fmt.Printf("ERROR: Couldn't write data to connection: %s", err)
	}

	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		input, err := rd.ReadString('\n')
		if err != nil {
			fmt.Printf("ERROR: Couldn't read from console: %s\n", err)
			continue
		}
		if input == "exit\n" {
			break
		}
		_, err = connection.Write([]byte(input))
		if err != nil {
			fmt.Printf("ERROR: Couldn't write data to connection: %s\n", err)
			continue
		}
	}
	_, err = connection.Write(close)
	if err != nil {
		fmt.Printf("ERROR: Couldn't write data to connection: %s", err)
	}
	connection.Close()
}
