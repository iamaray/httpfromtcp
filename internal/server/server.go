package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	Closed bool
}

func runConnection(s *Server, conn io.ReadCloser) {

}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.Closed {
			return 
		}
		if err != nil {
			return
		}
		go runConnection(s, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{ Closed: false }
	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.Closed = true
	return nil
}
