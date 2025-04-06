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

func (w *Writer) WriteStatusLine(statusCode int) error {
	if w.state != writerStateStatusLine {
		return errors.New("status write has been already wrote")
	}

	httpVersion := "HTTP/1.1"
	reasonPhrase := reasonPhrase(statusCode)

	statusLine := fmt.Sprintf("%s %d %s%s", httpVersion, statusCode, reasonPhrase, crfl)

	_, err := w.writer.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("error: writing status line err: %s", err)
	}

	w.state = writerStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateHeaders {
		return errors.New("to writer header first you must write the status line")
	}

	headers.Set("Connection", "close")

	_, hasContentLength := headers.Get("Content-Length")
	if !hasContentLength {
		return errors.New("Content-Length header is required")
	}

	_, hasContentType := headers.Get("Content-Type")
	if !hasContentType {
		headers.Set("Content-Type", "text/plain")
	}

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

	w.state = writerStateBody
	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, errors.New("to writer body first you must write the headers")
	}

	return w.writer.Write(body)
}
