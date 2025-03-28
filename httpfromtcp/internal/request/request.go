package request

import (
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
	lineBreak = "\r\n"

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

type Request struct {
	RequestLine RequestLine

	// isFullParsed holds the state of a request
	// while isFullParsed is false the request dint finish to parse yet
	isFullParsed bool

	rawRequest string
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := new(Request)
	var bytesReaded int
	var bytesParsed int

	buff := make([]byte, buffSize)
	for !request.isFullParsed {
		nFromReader, err := reader.Read(buff)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		buff = buff[:nFromReader]

		newBuff := make([]byte, 2*len(buff))
		if len(buff) == cap(buff) {
			_ = copy(newBuff, buff)
		}

		bytesReaded += nFromReader
		n, err := request.parse(buff)
		if err != nil {
			return nil, err
		}
		bytesParsed += n
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

// TODO: this implementation should parse data as
// the data come from and not store the all data it self and only parse when has full data in
func (r *Request) parseRequestLine(data []byte) (int, error) {
	dataStr := string(data)
	r.rawRequest += dataStr

	if !strings.Contains(r.rawRequest, lineBreak) {
		return 0, nil
	}

	requestPerLine := strings.Split(r.rawRequest, lineBreak)

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

	if err := r.validateTarget(target); err != nil {
		return 0, err
	}

	if err := r.validateHTTPVersion(fullHttpV, httpVSplited); err != nil {
		return 0, err
	}

	httpV := httpVSplited[1]

	r.RequestLine.HttpVersion = httpV
	r.RequestLine.RequestTarget = target
	r.RequestLine.Method = method

	return len(requestLine), nil
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

func (r *Request) validateTarget(target string) error {
	return nil
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
