package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

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

		// keep the connection open due some time
		// manage the state of how many conns opens do you have
		// be sure to not accept por conns whe you are full

		conn, err := s.tcpListener.Accept()
		if err != nil {
			log.Fatalf("Server - error on accept conn err: %s", err)
		}
		connID := newID()
		log.Printf("conn ID: %s - conn accepted", connID)

		go s.handleConn(conn, connID)
	}
}

func newID() string {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", newRand.Int63())
}

func (s *Server) handleConn(conn net.Conn, connID string) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("conn ID: %s - error o closing conn err: %s", err, connID)
		}
		log.Printf("conn ID: %s - conn closed", connID)
	}()

	request, err := request.ParseFromReader(conn)
	resp := response.NewWriter(conn)
	if err != nil {
		resp.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		resp.WriteHeaders(response.DefaultHeaders(len(body)))
		resp.WriteBody(body)
		return
	}

	s.handler(resp, request)
}
