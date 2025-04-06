package server

import (
	"fmt"

	"github.com/gpbPiazza/httpfromtcp/internal/headers"
	"github.com/gpbPiazza/httpfromtcp/internal/request"
	"github.com/gpbPiazza/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

func (he *HandlerError) Error() string {
	return ""
}

func (he *HandlerError) Write(w *response.Writer) error {
	if err := w.WriteStatusLine(he.StatusCode); err != nil {
		return err
	}

	headers := headers.New()
	headers.Add("Content-Length", fmt.Sprintf("%d", len(he.Message)))

	if err := w.WriteHeaders(headers); err != nil {
		return err
	}

	if _, err := w.WriteBody([]byte(he.Message)); err != nil {
		return err
	}

	return nil
}

type Handler func(w *response.Writer, req *request.Request) *HandlerError
