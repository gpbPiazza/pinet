package pihttp

import (
	"log"
	"strings"
)

func (s *Server) parseRequest(req []byte) Request {
	// TODO: IMPLEMENT PATH PARAM PARSE
	reqStr := string(req)

	reqByLineBreak := strings.Split(reqStr, lineBreak)

	reqLineIndex := 0
	entityBodyIndex := (len(reqByLineBreak) - 1)
	startHeadersIndex := reqLineIndex + 1
	endHeadersIndex := entityBodyIndex

	requestLine := parseRequestLine(reqByLineBreak[reqLineIndex])

	return Request{
		requestLine: requestLine,
		Headers:     parseHeaders(reqByLineBreak[startHeadersIndex:endHeadersIndex]),
		RawQuery:    requestLine.RawQuery(),
		EntityBody:  reqByLineBreak[entityBodyIndex],
	}
}

func parseRequestLine(requestLineStr string) requestLine {
	requestLineStr = strings.TrimSuffix(requestLineStr, lineBreak)
	requestLineSplit := strings.Split(requestLineStr, space)

	if len(requestLineSplit) < 3 {
		log.Fatal("uneexpected request line format: expected to have 3 elements inside of requestLineSplit")
	}

	return requestLine{
		Method:      requestLineSplit[0],
		Uri:         requestLineSplit[1],
		HttpVersion: requestLineSplit[2],
	}
}

func parseHeaders(headersStr []string) header {
	keyValSeparator := ": "
	header := make(header, len(headersStr))

	for _, headerStr := range headersStr {
		if headerStr == "" {
			continue
		}

		headerSplit := strings.Split(headerStr, keyValSeparator)
		if len(headerSplit) != 2 {
			log.Printf("header split in unexpected format -> headerStr: %s", headerStr)
			continue
		}

		key := headerSplit[0]
		val := headerSplit[1]

		val = strings.TrimPrefix(val, space)

		header[key] = append(header[key], val)
	}

	return header
}
