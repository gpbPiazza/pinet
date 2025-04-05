package main

import (
	"io"
	"log"

	"github.com/gpbPiazza/httpfromtcp/internal/request"
	"github.com/gpbPiazza/httpfromtcp/internal/response"
	"github.com/gpbPiazza/httpfromtcp/internal/server"
)

func main() {
	handler := func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:    "Your problem is not my problem",
			}
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "Woopsie, my bad",
			}
		}

		if req.RequestLine.RequestTarget == "/use-nvim" {
			if _, err := w.Write([]byte("All good, frfr")); err != nil {
				log.Print(err)
			}
			return nil
		}

		return &server.HandlerError{
			StatusCode: response.StatusNotFound,
			Message:    "the given request target is not found",
		}
	}

	server := server.New(server.WithHandler(handler))

	defer server.Close()

	server.Listen("42069")

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// <-sigChan
	// log.Println("Server gracefully stopped")
}
