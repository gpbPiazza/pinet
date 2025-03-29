package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	// carrieage return = \r
	// line feed = \n
	// carrieage return + line feed = CL
	crlf = "\r\n"

	// single space = SP
	space    = " "
	buffSize = 8

	httpVersionSuported = "HTTP/1.1"

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
	lineBreakBytes = []byte(crlf)
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

type Request struct {
	RequestLine RequestLine

	// isFullParsed holds the state of a request
	// while isFullParsed is false the request dint finish to parse yet
	isFullParsed bool
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := new(Request)
	var numBytesReaded int

	buff := make([]byte, buffSize)
	for !request.isFullParsed {
		if numBytesReaded >= len(buff) {
			newBuff := make([]byte, 2*len(buff))
			_ = copy(newBuff, buff)
			buff = newBuff
		}

		numBytesRead, err := reader.Read(buff[numBytesReaded:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		numBytesReaded += numBytesRead
		numBytesParsed, err := request.parse(buff[:numBytesReaded])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[numBytesParsed:])
		numBytesReaded -= numBytesParsed
	}

	return request, nil
}

func isAllCaps(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func (r *Request) parse(data []byte) (int, error) {
	if r.isFullParsed {
		return 0, errors.New("error: trying to parse data in a done state")
	}

	n, err := r.parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil
	}

	requestLineText := string(data[:idx])

	requestPerLine := strings.Split(requestLineText, crlf)

	requestLine := requestPerLine[0]

	requestLinePerSpace := strings.Split(requestLine, space)

	if len(requestLinePerSpace) != 3 {
		return 0, errors.New("request line has not 3 parts format")
	}

	method := requestLinePerSpace[0]
	target := requestLinePerSpace[1]
	fullHttpV := requestLinePerSpace[2]

	httpVSplited := strings.Split(fullHttpV, "/")

	if err := r.validateMethod(method); err != nil {
		return 0, err
	}

	if err := r.validateHTTPVersion(fullHttpV, httpVSplited); err != nil {
		return 0, err
	}

	httpV := httpVSplited[1]

	r.RequestLine.HttpVersion = httpV
	r.RequestLine.RequestTarget = target
	r.RequestLine.Method = method

	r.isFullParsed = true

	numBytesParsed := idx + 2 // +2 due CRLF

	return numBytesParsed, nil
}

func (r *Request) validateMethod(method string) error {
	if !isAllCaps(method) {
		return errors.New("request method malformed method is not in all captal letter")
	}

	for _, mappedM := range AllMethods {
		if mappedM == method {
			return nil
		}
	}

	return fmt.Errorf(
		"request method unsported - method got %s - see AllMethods variable to suported methods",
		method,
	)
}

func (r *Request) validateHTTPVersion(httpV string, httpVSplited []string) error {
	// TODO: add valiadtion to ensure version is digit . digit
	// now i am not looking if the version is a valid version

	if len(httpVSplited) != 2 {
		return errors.New("malformed http version expected <HTTP-NAME>/<digit>.<digit>")
	}

	if httpV != httpVersionSuported {
		return fmt.Errorf(
			"unsoported http version - the httpVersion is %s and only httpVersion suported is %s",
			httpV,
			httpVersionSuported,
		)
	}

	return nil
}
