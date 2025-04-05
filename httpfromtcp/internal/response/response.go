package response

import (
	"fmt"
	"io"
	"strings"

	"github.com/gpbPiazza/httpfromtcp/internal/headers"
)

const (
	crfl = "\r\n"
)

func WriteStatusLine(w io.Writer, statusCode int) error {
	httpVersion := "HTTP/1.1"
	reasonPhrase := reasonPhrase(statusCode)

	statusLine := fmt.Sprintf("%s %d %s%s", httpVersion, statusCode, reasonPhrase, crfl)

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("error: writing status line err: %s", err)
	}

	return nil
}

func DefaultHeaders(contentLen int) headers.Headers {
	headers := headers.New()

	headers.Add("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Add("Connection", "close")
	headers.Add("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
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

	_, err := w.Write([]byte(fieldLines.String()))
	if err != nil {
		return fmt.Errorf("error on write headers err: %s", err)
	}

	return nil
}

type Writer struct {
}

func (w *Writer) WriteStatusLine(statusCode int) error
func (w *Writer) WriteHeaders(headers headers.Headers) error
func (w *Writer) WriteBody(p []byte) (int, error)
