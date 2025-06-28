package main

import (
	"fmt"
	"github.com/lucoand/httpfromtcp/internal/request"
	// "io"
	"log"
	"net"
	// "strings"
)

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	ch := make(chan string)
//
// 	go func() {
// 		defer f.Close()
// 		defer close(ch)
// 		content := make([]byte, 8)
// 		line := ""
// 		for {
// 			n, err := f.Read(content)
// 			if err == io.EOF {
// 				break
// 			}
// 			if err != nil {
// 				log.Fatalf("Error reading from file: %s", err)
// 			}
//
// 			contentString := string(content[:n])
// 			contentParts := strings.Split(contentString, "\n")
// 			line += contentParts[0]
// 			if len(contentParts) > 1 {
// 				for _, part := range contentParts[1:] {
// 					ch <- line
// 					line = part
// 				}
// 			}
// 		}
// 		if line != "" {
// 			ch <- line
// 		}
// 	}()
//
// 	return ch
// }

func main() {
	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("ERROR: couldn't open net listener: %s", err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: couldn't connect: %s", err)
		}
		fmt.Printf("Connection established!\n")
		// ch := getLinesChannel(connection)
		// for line := range ch {
		// 	fmt.Printf("%s\n", line)
		// }
		request, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatalf("ERROR: could parse request: %s", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Printf("Connection has been closed.\n")
	}
}
