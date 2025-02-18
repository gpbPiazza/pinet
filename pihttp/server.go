package pihttp

import (
	"log"
	"net"
	"strings"
)

// -> https://pkg.go.dev/net#pkg-examples
// -> https://itsfuad.medium.com/understanding-http-at-a-low-level-a-developers-guide-with-c-213728d6c41d

// Create a struct that is capeable to stablish a connection using TCP/IP protocol from net package. DONE
// This struct shoul be able to handle clients to connect in they open socket IP connection and response the request client  DONE

// Make this server accept HTTP requests and read thoose requests
// Make the server be able to parse queryParams
// Make the server be able to response with HTTP response format

// Make the server be able to parse bodyParams
// Make this server be able differ from POST, GET, PUT, PATCH, DELETE methods

const (
	DefaultWriteBufferSize = 4096
	DefaultReadBufferSize  = 4096
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
			log.Printf("Server - error on listener close err: %s", err)
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

	s.writeResp(conn)
}

// reading HTTP RFC -> https://www.rfc-editor.org/rfc/rfc1945.html#section-5
// SECTION THAT DEFINE A HTTP REQUEST FORMAT
func (s *Server) parseRequest() {
	// "GET / HTTP/1.1\r\nHost: localhost:8080\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n"

	// HTTP REQUEST FORMAT
	// Request-Line = Method SP Request-URI SP HTTP-Version CRLF

	// if a higher version request is received, the
	//  proxy/gateway must either downgrade the request version or respond
	//  with an error.

}

// reading HTTP RFC -> https://www.rfc-editor.org/rfc/rfc1945.html#section-6
// SECTION THAT DEFINE A HTTP RESPONSE FORMAT
func (s *Server) writeResp(conn net.Conn) {
	//  After receiving and interpreting a request message, a server responds
	//  in the form of an HTTP response message.
	//  Response        = Simple-Response | Full-Response
	//  Simple-Response = [ Entity-Body ]

	// Simple response should be only from HTTP versions <= 0.9.
	// If a client sends an HTTP/1.0 Full-Request and the server response with
	// with a Status-Line the client should assume that is a Simple-response.

	// Simple response format:
	// Just entity body -> see definition of entity body

	// Full response format:

	// entity body definition
	// Entity Body

	//  The entity body (if any) sent with an HTTP request or response is in
	//  a format and encoding defined by the Entity-Header fields.
	//      Entity-Body    = *OCTET

	// if a request has a body that means that in the request we have:
	// 1. the http request metthod allows.
	// 2. We have Content-Length header field

	//  For response messages,
	// All responses dependent on request method and response code.
	// All responses to the HEAD request method must not include a body.
	// All status code 1xx (informational), 204 (no content), and 304 (not modified) responses must not include a body.
	// All other responses must include body or a Content-Length header field defined with a value of zero (0).

	strBuilder := new(strings.Builder)

	strBuilder.WriteString("HTTP/1.1 200 OK\r\n")

	res := strBuilder.String()

	// std::stringstream response;
	// response << "HTTP/1.1 200 OK\r\n";
	// response << "Content-Type: text/html\r\n";
	// response << "Content-Length: 46\r\n"; // Length of the HTML content
	// response << "\r\n";
	// response << "<html><body><h1>Hello, World!</h1></body></html>";
	nWritten, err := conn.Write([]byte(res))
	if err != nil {
		log.Fatalf("Server - error on Write client conn err: %s", err)
	}
	log.Printf("Server - number of bytes successfully written into conn: %d", nWritten)
}
