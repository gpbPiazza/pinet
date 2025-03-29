package httpfromtcp

import (
	"fmt"
	"log"
	"net"

	"github.com/gpbPiazza/httpfromtcp/internal/request"
)

type Server struct {
}

func NewServer(opts ...Option) *Server {
	option := options{}

	for _, opt := range opts {
		opt.apply(&option)
	}

	s := &Server{}

	return s
}

func (s *Server) Listen(address string) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", address))
	if err != nil {
		log.Fatalf("Server - error on create listener conn err: %s", err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("Server - error on listener close err: %s", err)
		}
	}()

	log.Printf("starting listener at port: %d", 42069)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Server - error on accept conn err: %s", err)
		}
		log.Print("conn accepted")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("error on parse request err: %s", err)
		}

		fmt.Print("Request line:\n")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
	}
}
