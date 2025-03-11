package pihttp

import (
	"bytes"
	"log"
)

func (s *Server) parseRequest(req []byte) Request {
	reqByLineBreak := bytes.Split(req, lineBreakBytes)

	reqLineIndex := 0
	entityBodyIndex := len(reqByLineBreak) - 1
	startHeadersIndex := reqLineIndex + 1
	endHeadersIndex := entityBodyIndex

	requestLine := parseRequestLine(reqByLineBreak[reqLineIndex])

	headers := parseHeaders(reqByLineBreak[startHeadersIndex:endHeadersIndex])

	return Request{
		requestLine: requestLine,
		Headers:     headers,
		RawQuery:    requestLine.RawQuery(),
		EntityBody:  reqByLineBreak[entityBodyIndex],
	}
}

func parseRequestLine(requestLineBytes []byte) requestLine {
	// Trim the line break from the end of the request line
	requestLineBytes = bytes.TrimSuffix(requestLineBytes, lineBreakBytes)
	requestLineSplit := bytes.Split(requestLineBytes, spaceBytes)

	if len(requestLineSplit) < 3 {
		log.Fatal("unexpected request line format: expected to have 3 elements inside of requestLineSplit")
	}

	uri := string(requestLineSplit[1])

	path, _, ok := bytes.Cut(requestLineSplit[1], queryDelimiterBytes)
	if !ok {
		// TODO: validate if a request without queryParams will return not OK here
		// I think that will because ? is optional
		log.Printf("not OK in cut path from URI")
	}

	return requestLine{
		Method:      string(requestLineSplit[0]),
		Uri:         uri,
		Path:        string(path),
		HttpVersion: string(requestLineSplit[2]),
	}
}

func parseHeaders(headersBytes [][]byte) Header {
	keyValSeparator := []byte(": ")
	header := make(Header, len(headersBytes))

	for _, headerBytes := range headersBytes {
		if len(headerBytes) == 0 {
			continue
		}

		headerSplit := bytes.SplitN(headerBytes, keyValSeparator, 2)
		if len(headerSplit) != 2 {
			log.Printf("header split in unexpected format -> headerBytes: %s", headerBytes)
			continue
		}

		key := string(headerSplit[0])
		val := string(bytes.TrimPrefix(headerSplit[1], spaceBytes))

		header[key] = append(header[key], val)
	}

	return header
}
