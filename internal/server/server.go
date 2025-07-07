package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/lucoand/httpfromtcp/internal/request"
	"github.com/lucoand/httpfromtcp/internal/response"
)

type Server struct {
	IsClosed *atomic.Bool
	Listener net.Listener
	Handler  Handler
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode string
	Message    string
}

func (s *Server) listen() {
	for !s.IsClosed.Load() {
		conn, err := s.Listener.Accept()
		if err != nil && !s.IsClosed.Load() {
			continue
		}
		if !s.IsClosed.Load() {
			fmt.Println("Handling request")
			go s.handle(conn)
		}
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Parsing request")
	req, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}
	fmt.Println("Request parsed")

	buf := &bytes.Buffer{}
	fmt.Println("Buffer created")
	handlerError := s.Handler(buf, req)
	fmt.Println("Handler called")
	if handlerError != nil {
		writeHandlerError(conn, handlerError)
		return
	}
	h := response.GetDefaultHeaders(buf.Len())
	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		return
	}
	err = response.WriteHeaders(conn, h)
	if err != nil {
		return
	}

	io.Copy(conn, buf)

}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.IsClosed.Store(true)
	return nil
}

func writeHandlerError(w io.Writer, h *HandlerError) error {
	errorString := "HTTP/1.1 " + h.StatusCode + " " + h.Message + "\r\n"
	errorBytes := []byte(errorString)
	n, err := w.Write(errorBytes)
	return response.WriteErrorHelper(err, n, errorBytes)
}

func Serve(port int, h Handler) (*Server, error) {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	var isClosed atomic.Bool
	isClosed.Store(false)
	s := Server{
		Listener: listener,
		IsClosed: &isClosed,
		Handler:  h,
	}
	fmt.Println("Handler attached")
	go s.listen()
	return &s, nil
}
