package server

import (
	"github.com/gpbPiazza/httpfromtcp/internal/request"
	"github.com/gpbPiazza/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
