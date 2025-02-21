package pihttp

import (
	"fmt"
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

	req := s.parseRequest(buffer[:nReaded])

	log.Print("Request Line Method: ", req.method)
	log.Print("Request Line URI: ", req.uri)
	log.Print("Request Line HTTP Version: ", req.httpVersion)

	s.writeResp(conn)
}

type requestLine struct {
	method      string
	uri         string
	httpVersion string
}

func parseRequestLine(requestLineStr string) requestLine {
	requestLineStr = strings.TrimSuffix(requestLineStr, lineBreak)
	requestLineSplit := strings.Split(requestLineStr, space)

	if len(requestLineSplit) < 3 {
		log.Fatal("uneexpected request line format: expected to have 3 elements inside of requestLineSplit")
	}

	return requestLine{
		method:      requestLineSplit[0],
		uri:         requestLineSplit[1],
		httpVersion: requestLineSplit[2],
	}
}

// key value pair
// headers can have repetitive key names, if so,
// this values will me concatenated into the same key
type header map[string][]string

// User-Agent: Go-http-client/1.1
// Accept-Encoding: gzip
// Times-Do-RJ: fla, flu, vasco
// Times-Do-RJ: botafogo
// Accept-Encoding: gremio
func parseHeaders(headersStr []string) header {
	keyValSeparator := ":"
	header := make(header, len(headersStr))

	for _, headerStr := range headersStr {
		headerSplit := strings.Split(headerStr, keyValSeparator)
		key := headerSplit[0]
		val := headerSplit[1]

		// manyVals := strings.Split(val, ",")

		val = strings.TrimPrefix(val, space)

		header[key] = append(header[key], val)
	}

	return header
}

type request struct {
	requestLine
	header header
}

// reading HTTP RFC -> https://www.rfc-editor.org/rfc/rfc1945.html#section-5
// SECTION THAT DEFINE A HTTP REQUEST FORMAT
func (s *Server) parseRequest(req []byte) request {
	// "GET / HTTP/1.1\r\n
	// Host: localhost:8080\r\n
	// User-Agent: Go-http-client/1.1\r\n
	// Accept-Encoding: gzip\r\n
	// \r\n"

	reqStr := string(req)

	reqByLineBreak := strings.Split(reqStr, lineBreak)

	requestLineStr := reqByLineBreak[0]

	headersAndEntityBody := reqByLineBreak[0:]

	var headers []string
	// var entityBody string
	for _, reqVal := range headersAndEntityBody {
		if reqVal != lineBreak {
			headers = append(headers, strings.TrimSuffix(reqVal, lineBreak))
		}
		// TODO não sei como que vem o body ainda...
	}

	requestLine := parseRequestLine(requestLineStr)
	header := parseHeaders(headers)

	// HTTP REQUEST FORMAT
	// Request-Line = Method SP Request-URI SP HTTP-Version CRLF

	// if a higher version request is received, the
	//  proxy/gateway must either downgrade the request version or respond
	//  with an error.
	return request{
		requestLine: requestLine,
		header:      header,
	}
}

type response struct {
	statusLine     string
	generalHeader  string
	responseHeader string
	entityHeader   string
	entityBody     string
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
