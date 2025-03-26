package httpfromtcp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

	log.Printf("starting liestener at port: %d", 42069)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Server - error on accept conn err: %s", err)
		}
		log.Print("conn accepted")

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("%s\n", line)
		}

	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)

	go func() {
		defer func() {
			log.Print("close chann and conn")
			if err := f.Close(); err != nil {
				log.Printf("err while reading closing conn err: %s", err)
			}
			close(linesChan)
		}()

		var err error
		var parts []string
		for !errors.Is(err, io.EOF) {
			data := make([]byte, 8)
			_, err = f.Read(data)
			if errors.Is(err, io.EOF) {
				line := strings.Join(parts, "")
				linesChan <- line

				break
			}

			if err != nil {
				log.Printf("err while reading from conn err: %s", err.Error())
			}

			nLine := "\n"
			dataStr := string(data)

			lineSegments := strings.Split(dataStr, nLine)

			hasLineCut := len(lineSegments) > 1
			if !hasLineCut {
				parts = append(parts, dataStr)
				continue
			}

			parts = append(parts, lineSegments[0])
			line := strings.Join(parts, "")

			linesChan <- line

			parts = nil
			parts = append(parts, lineSegments[1:]...)
		}
	}()

	return linesChan
}
