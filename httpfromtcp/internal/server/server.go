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
}

func New(opts ...Option) *Server {
	option := options{}

	for _, opt := range opts {
		opt.apply(&option)
	}

	closed := &atomic.Bool{}
	closed.Store(false)

	s := &Server{
		isClosed: closed,
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

	if err != nil {
		log.Printf("error on parse request err: %s", err)
	}

	fmt.Print("Request line:\n")
	fmt.Printf("- Method: %s\n", request.RequestLine.Method)
	fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)

	fmt.Print("Headers:\n")
	for key, val := range request.Headers {
		fmt.Printf("- %s: %s\n", key, val)
	}
	fmt.Print("Body:\n")
	fmt.Printf("%s\n", string(request.Body))

	respHeaders := response.DefaultHeaders(len(request.Body))
	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		log.Print(err)
	}
	if err := response.WriteHeaders(conn, respHeaders); err != nil {
		log.Print(err)
	}
}
