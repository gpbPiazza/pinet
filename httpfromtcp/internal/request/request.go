package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/gpbPiazza/httpfromtcp/internal/headers"
)

const (
	// carrieage return + line feed
	crlf = "\r\n"
	// single space = SP
	space    = " "
	buffSize = 8

	httpVSuported = "HTTP/1.1"

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
	Headers     headers.Headers
	state       requestState
}

type requestState int

const (
	requestStateInitialized = iota
	requestStateParsingHeaders
	requestStateCompled
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func ParseFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state:   requestStateInitialized,
		Headers: headers.New(),
	}

	var numBytesReaded int

	buff := make([]byte, buffSize)
	for !request.isFullParsed() {
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

		newBuff := make([]byte, len(buff))
		buffDst := buff[numBytesParsed:]
		copy(newBuff, buffDst)
		buff = newBuff
		numBytesReaded -= numBytesParsed
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != requestStateCompled {
		toBeParsed := data[totalBytesParsed:]

		numBytesParsed, err := r.parseSingle(toBeParsed)
		totalBytesParsed += numBytesParsed

		if err != nil {
			return 0, err
		}

		if totalBytesParsed > len(toBeParsed) {
			return totalBytesParsed, nil
		}

		if numBytesParsed == 0 {
			return totalBytesParsed, nil
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		n, err := r.parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateCompled
		}
		return n, nil
	case requestStateCompled:
		return 0, errors.New("error: trying to parse data in a done state")
	default:
		return 0, errors.New("unknow request state")
	}
}

// parseRequestLine will keep track of data until has all requestLine in data then will parse requestLine
// parseRequestLine will set RequestLine values into request.
func (r *Request) parseRequestLine(data []byte) (int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil
	}

	requestText := string(data[:idx])
	requestPerLine := strings.Split(requestText, crlf)
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

	numBytesParsed := len(requestLine) + 2 // +2 due CRLF

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

func (r *Request) validateHTTPVersion(httpV string, httpVSplited []string) error {
	if len(httpVSplited) != 2 {
		return errors.New("malformed http version expected <HTTP-NAME>/<digit>.<digit>")
	}

	if httpV != httpVSuported {
		return fmt.Errorf(
			"unsoported http version - the httpVersion is %s and only httpVersion suported is %s",
			httpV,
			httpVSuported,
		)
	}

	return nil
}

func (r *Request) isFullParsed() bool {
	return r.state == requestStateCompled
}
