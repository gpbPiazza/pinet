package pihttp

import (
	"log"
	"net"
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
	// carrieage return + line feed = CL
	lineBreak = "\r\n"

	// single space = SP
	space = " "

	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

var (
	lineBreakBytes = []byte(lineBreak)
	spaceBytes     = []byte(space)

	AllMethods = []string{
		MethodGet,
		MethodHead,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
		MethodConnect,
		MethodOptions,
		MethodTrace,
	}
)

func NewServer() *Server {
	s := &Server{}

	s.routes = make(map[string]map[string]Handler)
	for _, method := range AllMethods {
		s.routes[method] = make(map[string]Handler)
	}

	return s
}

type Server struct {
	routes map[string]map[string]Handler
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

	handlerByRoute, ok := s.routes[req.Method]
	if !ok {
		// TODO: implement not found route resp
		log.Fatalf("request method not registered in routes map. Received request method: %s", req.Method)
	}

	handler, ok := handlerByRoute[req.Path]
	if !ok {
		// TODO: implement not found route resp
		log.Fatalf("request Path not registered in routes map. Received request Path: %s", req.Path)
	}

	resp := newResp(req)
	if err := handler(req, resp); err != nil {
		// TODO: implement write response in error cases
		log.Printf("errro from client handler err: %s", err)
	}

	s.writeResp(conn, *resp)
}

func newResp(req Request) *Response {
	resp := new(Response)
	resp.Header = make(Header)

	resp.httpVersion = req.HttpVersion

	return resp
}

func (s *Server) writeResp(conn net.Conn, resp Response) {
	nWritten, err := conn.Write(resp.encode())
	if err != nil {
		log.Fatalf("Server - error on Write client conn err: %s", err)
	}
	log.Printf("Server - number of bytes successfully written into conn: %d", nWritten)
}

type Handler func(req Request, resp *Response) error

func (s *Server) HandleFunc(method, path string, handler Handler) {
	s.routes[method][path] = handler
}
