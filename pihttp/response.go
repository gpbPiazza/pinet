package pihttp

import (
	"bytes"
	"fmt"
	"log"
)

type Response struct {
	Header      Header
	StatusCode  int
	Body        []byte
	httpVersion string
	// statusLine     string
	// generalHeader  string
	// responseHeader string
	// entityHeader   string
	// entityBody     string
}

// encode will encode Response struct into valid http response format to []byte.
func (r *Response) encode() []byte {
	// ITS IS MUST PASS:
	// REQUEST LINE RESPONSE HTTP VERSION STATUS CODE STATUS MESSAGE
	// CONTENT TYPE HEADER
	// CONTENT LENGTH
	// BREAK LINE
	// ENTITYBody

	respBuff := bytes.NewBuffer(nil)

	responseLine := fmt.Sprintf("%s%s%d%s%s%s", r.httpVersion, space, r.StatusCode, space, statusText(r.StatusCode), lineBreak)
	_, err := respBuff.WriteString(responseLine)
	if err != nil {
		// TODO: i dont know if this should be a fatal from my lib
		// probably this should be error wrapped to the client treat in some error response middleware
		log.Fatalf("error on write respBuff RESPONSE LINE err: %s", err)
	}

	hasContenTypeHeader := false
	hasContenLengthHeader := false
	for headerKey, headerVal := range r.Header {
		if headerKey == "Cotent-Type" {
			hasContenTypeHeader = true
		}
		if headerKey == "Content-Length" {
			hasContenLengthHeader = true
		}

		for _, val := range headerVal {
			_, err := respBuff.WriteString(fmt.Sprintf("%s: %s%s", headerKey, val, lineBreak))
			if err != nil {
				// TODO: i dont know if this should be a fatal from my lib
				// probably this should be error wrapped to the client treat in some error response middleware
				log.Fatalf("error on write respBuff HEADERS err: %s", err)
			}
		}
	}

	if !hasContenLengthHeader {
		_, err := respBuff.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(r.Body)))
		if err != nil {
			// TODO: i dont know if this should be a fatal from my lib
			// probably this should be error wrapped to the client treat in some error response middleware
			log.Fatalf("error on write respBuff HEADER CONTENT LENGTH err: %s", err)
		}
	}

	if !hasContenTypeHeader {
		// TODO identify type from resp
		_, err := respBuff.WriteString("Content-Type: text/html\r\n")
		if err != nil {
			// TODO: i dont know if this should be a fatal from my lib
			// probably this should be error wrapped to the client treat in some error response middleware
			log.Fatalf("error on write respBuff HEADER CONTENT TYPE err: %s", err)
		}
	}

	_, err = respBuff.Write(lineBreakBytes)
	if err != nil {
		// TODO: i dont know if this should be a fatal from my lib
		// probably this should be error wrapped to the client treat in some error response middleware
		log.Fatalf("error on write respBuff BreakLine Between end Headers and start entityBody err: %s", err)
	}

	_, err = respBuff.Write(r.Body)
	if err != nil {
		// TODO: i dont know if this should be a fatal from my lib
		// probably this should be error wrapped to the client treat in some error response middleware
		log.Fatalf("error on write respBuff ENTITYBODY err: %s", err)
	}

	return respBuff.Bytes()
}
