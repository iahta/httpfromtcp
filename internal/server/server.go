package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/iahta/httpfromtcp/internal/request"
	"github.com/iahta/httpfromtcp/internal/response"
)

type Server struct {
	Listener net.Listener
	Handler  Handler
	Closed   atomic.Bool
}

func Serve(port int, h Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("unable to serve listener: %v", err)
	}
	server := &Server{
		Listener: l,
		Handler:  h,
	}

	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.Closed.Store(true)
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.Closed.Load() {
				break
			}
			fmt.Printf("error making connection: %v", err)
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.Handler(w, req)
}
