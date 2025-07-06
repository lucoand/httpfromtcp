package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/lucoand/httpfromtcp/internal/response"
)

type Server struct {
	IsClosed *atomic.Bool
	Listener net.Listener
}

func (s *Server) listen() {
	for !s.IsClosed.Load() {
		conn, err := s.Listener.Accept()
		if err != nil && !s.IsClosed.Load() {
			continue
		}
		if !s.IsClosed.Load() {
			go s.handle(conn)
		}
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		return
	}
	h := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, h)
	if err != nil {
		return
	}
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.IsClosed.Store(true)
	return nil
}

func Serve(port int) (*Server, error) {
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
	}
	go s.listen()
	return &s, nil
}
