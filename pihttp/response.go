package pihttp

import (
	"fmt"
	"strings"
)

type Response struct {
	Header     Header
	StatusCode int
	Body       any
	// statusLine     string
	// generalHeader  string
	// responseHeader string
	// entityHeader   string
	// entityBody     string
}

// encode will parse Response struct into valid http response format.
func (r *Response) encode() []byte {
	// ITS IS MUST PASS:
	// REQUEST LINE RESPONSE HTTP VERSION STATUS CODE STATUS MESSAGE
	// CONTENT TYPE HEADER
	// CONTENT LENGTH
	// BREAK LINE
	// ENTITYBody

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

	return []byte(res)
}
