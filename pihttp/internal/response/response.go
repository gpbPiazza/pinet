package response

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gpbPiazza/httpfromtcp/internal/headers"
)

const (
	crfl = "\r\n"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
	writerStateTrailers
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  writerStateStatusLine,
	}
}

func DefaultHeaders(bodyLen int) headers.Headers {
	h := headers.New()

	h.Override("Connection", "close")
	h.Override("Content-Type", "text/plain")
	h.Override("Content-Length", fmt.Sprintf("%d", bodyLen))

	return h
}

func (w *Writer) WriteStatusLine(statusCode int) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.state)
	}
	defer func() { w.state = writerStateHeaders }()

	httpVersion := "HTTP/1.1"
	reasonPhrase := reasonPhrase(statusCode)

	statusLine := fmt.Sprintf("%s %d %s%s", httpVersion, statusCode, reasonPhrase, crfl)

	_, err := w.writer.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("error: writing status line err: %s", err)
	}

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.state)
	}
	defer func() { w.state = writerStateBody }()

	fieldLines := new(strings.Builder)
	for key, val := range headers {
		fiedlLine := fmt.Sprintf("%s: %s%s", key, val, crfl)

		if _, err := fieldLines.WriteString(fiedlLine); err != nil {
			return fmt.Errorf("error to write field line err: %s", err)
		}
	}

	if _, err := fieldLines.WriteString(crfl); err != nil {
		return fmt.Errorf("error to write field line err: %s", err)
	}

	_, err := w.writer.Write([]byte(fieldLines.String()))
	if err != nil {
		return fmt.Errorf("error on write headers err: %s", err)
	}

	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, errors.New("to writer body first you must write the headers")
	}

	defer func() { w.state = writerStateTrailers }()

	return w.writer.Write(body)
}

// WriteChunkedBody will write into conn the body in following format below:
//
// HTTP/1.1 200 OK
// Content-Type: text/plain
// Transfer-Encoding: chunked
//
// <n>/r/n
// <data of length n>/r/n
// <n>/r/n
// <data of length n>/r/n
// <n>/r/n
// <data of length n>/r/n
// <n>/r/n
// <data of length n>/r/n
// ... repeat ...
// 0\r\n
// \r\n
//
// Where n is the len of the bytes content and the next line is the content.
//
// WriteChunkedBody will write a new line in the chunked body to each call.
//
// To finish writing into chunked body you must call WriteChunkedBodyDone.
func (w *Writer) WriteChunkedBody(chunk []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}
	chunkSize := len(chunk)

	numberTotal := 0
	numberBytes, err := fmt.Fprintf(w.writer, "%x%s", chunkSize, crfl)
	if err != nil {
		return numberTotal, err
	}
	numberTotal += numberBytes

	numberBytes, err = w.writer.Write(chunk)
	if err != nil {
		return numberTotal, err
	}
	numberTotal += numberBytes

	numberBytes, err = w.writer.Write([]byte(crfl))
	if err != nil {
		return numberTotal, err
	}
	numberTotal += numberBytes
	return numberTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.state)
	}

	defer func() { w.state = writerStateTrailers }()

	return w.writer.Write([]byte(fmt.Sprintf("%d%s", 0, crfl)))
}

func (w *Writer) WriteTrailers(header headers.Headers) error {
	if w.state != writerStateTrailers {
		return fmt.Errorf("cannot write trailer body in state %d", w.state)
	}

	trailers, ok := header.Get("Trailer")
	if !ok {
		return errors.New("no trailer header key")
	}

	trailerNames := strings.Split(trailers, headers.ValSeparator)

	for _, tName := range trailerNames {
		tVal, ok := header.Get(tName)
		if !ok {
			return errors.New("registered trailer name not present into header")
		}

		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s%s", tName, tVal, crfl)))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte(crfl))
	if err != nil {
		return err
	}

	return nil
}
