package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/gpbPiazza/httpfromtcp/internal/request"
	"github.com/gpbPiazza/httpfromtcp/internal/response"
)

type Server struct {
	tcpListener net.Listener
	isClosed    *atomic.Bool

	handler Handler
}

func New(opts ...Option) *Server {
	option := options{
		handler: nil,
	}

	for _, opt := range opts {
		opt.apply(&option)
	}

	closed := new(atomic.Bool)
	closed.Store(false)

	s := &Server{
		isClosed: closed,
		handler:  option.handler,
	}

	return s
}

func (s *Server) Close() error {
	defer s.isClosed.Store(true)

	return s.tcpListener.Close()
}

func (s *Server) Listen(address string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", address))
	if err != nil {
		log.Fatalf("Server - error on create listener conn err: %s", err)
	}
	s.tcpListener = listener

	log.Printf("starting listener at port: %d", 42069)

	for {
		if s.isClosed.Load() {
			break
		}

		conn, err := s.tcpListener.Accept()
		if err != nil {
			log.Fatalf("Server - error on accept conn err: %s", err)
		}
		log.Print("conn accepted")

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error o closing conn err: %s", err)
		}
	}()

	request, err := request.ParseFromReader(conn)
	resp := response.NewWriter(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		if err := hErr.Write(resp); err != nil {
			log.Printf("error to write into conn handler err: %s", err)
		}
		return
	}

	hErr := s.handler(resp, request)
	if hErr != nil {
		if err := hErr.Write(resp); err != nil {
			log.Printf("error to write into conn handler err: %s", err)
		}
		return
	}
}
