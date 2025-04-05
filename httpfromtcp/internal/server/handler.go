package server

import (
	"io"

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

func (he *HandlerError) Write(w io.Writer) error {
	if err := response.WriteStatusLine(w, he.StatusCode); err != nil {
		return err
	}

	respHeaders := response.DefaultHeaders(len(he.Message))

	if err := response.WriteHeaders(w, respHeaders); err != nil {
		return err
	}

	if _, err := w.Write([]byte(he.Message)); err != nil {
		return err
	}

	return nil
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
