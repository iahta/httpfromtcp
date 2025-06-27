package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	Closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("unable to serve listener: %v", err)
	}
	server := &Server{
		Listener: l,
	}

	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.Closed.Store(true)

	if s.Closed.CompareAndSwap(false, true) {
		fmt.Printf("Successfully closed server")
	}
	return nil
}

func (s *Server) listen() {
	for !s.Closed.Load() {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.Closed.Load() {
				break
			}
			fmt.Printf("error making connection: %v", err)
		}
		if conn != nil {
			go s.handle(conn)
		}
	}
}

func (s *Server) handle(conn net.Conn) {
	resp := []byte("HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n\r\n" +
		"Hello World!",
	)

	defer conn.Close()
	_, err := conn.Write(resp)
	if err != nil {
		fmt.Printf("error writing response: %v", err)
	}

}
