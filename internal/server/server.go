package server

import (
	"bytes"
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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	hErr := s.Handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return

}
