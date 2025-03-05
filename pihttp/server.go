package pihttp

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// -> https://pkg.go.dev/net#pkg-examples
// -> https://itsfuad.medium.com/understanding-http-at-a-low-level-a-developers-guide-with-c-213728d6c41d

// Make the server be able to response with HTTP response format
// Make this server be able differ from POST, GET, PUT, PATCH, DELETE methods

const (
	DefaultWriteBufferSize = 4096
	DefaultReadBufferSize  = 4096

	// carrieage return = \r
	// line feed = \n
	lineBreak = "\r\n"
	// single space = SP
	space = " "
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Server - error on create listener conn err: %s", err)
	}

	log.Printf("Server - CONN local Addr network: %s", listener.Addr().Network())
	log.Printf("Server - CONN local literal Addrs: %s", listener.Addr().String())

	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("Server - error on listener close err: %s", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Server - error on accept conn err: %s", err)
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Server - error on close conn err: %s", err)
		}
	}()
	log.Printf("Client CONN - Addr network: %s", conn.LocalAddr().Network())
	log.Printf("Client CONN - literal Addrs: %s", conn.LocalAddr().String())

	buffer := make([]byte, DefaultReadBufferSize)
	nReaded, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("error on Read client conn err: %s", err)
	}
	log.Printf("Client request - number of bytes read from the conn: %d", nReaded)
	log.Printf("Client request - %s", string(buffer[:nReaded]))

	req := s.parseRequest(buffer[:nReaded])

	log.Print("Request Line Method: ", req.Method)
	log.Print("Request Line URI: ", req.Uri)
	log.Print("Request Line HTTP Version: ", req.HttpVersion)

	s.writeResp(conn)
}

func (s *Server) writeResp(conn net.Conn) {
	strBuilder := new(strings.Builder)

	entityBody := `
		<html>
			<body>
				<h1>Hello, World!</h1>
			</body>
		</html>
	`

	strBuilder.WriteString("HTTP/1.1 200 OK\r\n")
	strBuilder.WriteString("Content-Type: text/html\r\n")
	strBuilder.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len([]byte(entityBody))))
	strBuilder.WriteString("\r\n")
	strBuilder.WriteString(entityBody)

	res := strBuilder.String()

	nWritten, err := conn.Write([]byte(res))
	if err != nil {
		log.Fatalf("Server - error on Write client conn err: %s", err)
	}
	log.Printf("Server - number of bytes successfully written into conn: %d", nWritten)
}

type Handler func(req Request, resp *Response) error

func (s *Server) HandleFunc(method, path string, handler Handler) {
	// map[method+path]handler

	// handleConn
	// faz parse
	// pega method+path
	// executam handler
	// parse clientResponse
}
