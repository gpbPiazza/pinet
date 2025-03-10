package pihttp

import (
	"log"
	"strings"
)

func (s *Server) parseRequest(req []byte) Request {
	// TODO: IMPLEMENT PATH PARAM PARSE

	// when i do this i parse byts of 8bit patterns to
	// decode to  UTF-8-encoded text, there is some way to work with req
	// without parsing to UTF-8?
	// OK there is, for all method using strings.Split, strings we can use
	// bytes package
	// for now I am ok with this unnecerring string parsing
	// very go TODO to make, register some banch marking with stringds pakcage
	// implement with bytes package and see if there is any significan change dont do this kind
	// of encoding in middle of parsing request.
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
		// TODO: FIX THIS -> IN THE server implementation of json unmarshall he got
		// 2025/03/10 09:02:01 errro from client handler err: json: cannot unmarshal string into Go value of type main.TextTimeBody
		EntityBody: []byte(reqByLineBreak[entityBodyIndex]),
	}
}

func parseRequestLine(requestLineStr string) requestLine {
	requestLineStr = strings.TrimSuffix(requestLineStr, lineBreak)
	requestLineSplit := strings.Split(requestLineStr, space)

	if len(requestLineSplit) < 3 {
		log.Fatal("uneexpected request line format: expected to have 3 elements inside of requestLineSplit")
	}

	uri := requestLineSplit[1]

	path, _, ok := strings.Cut(uri, queryDelimiter)
	if !ok {
		// TODO: validate if a request without queryParams will return not OK here
		// I think that will because ? is optional
		log.Printf("not OK in cut path from URI")
	}

	return requestLine{
		Method:      requestLineSplit[0],
		Uri:         uri,
		Path:        path,
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
