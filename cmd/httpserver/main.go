package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lucoand/httpfromtcp/internal/request"
	"github.com/lucoand/httpfromtcp/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	fmt.Println("Handler entered")
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: "400",
			Message:    "Your problem is not my problem\n",
		}
	}
	if target == "/myproblem" {
		return &server.HandlerError{
			StatusCode: "500",
			Message:    "Woopsie, my bad\n",
		}
	}
	bodyString := "All good, frfr\n"
	bodyBytes := []byte(bodyString)
	w.Write(bodyBytes)
	return nil
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
