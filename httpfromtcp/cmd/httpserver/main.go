package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gpbPiazza/httpfromtcp/internal/headers"
	"github.com/gpbPiazza/httpfromtcp/internal/request"
	"github.com/gpbPiazza/httpfromtcp/internal/response"
	"github.com/gpbPiazza/httpfromtcp/internal/server"
)

func main() {
	handler := func(w *response.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return handleBadRequest(w, req)
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			return handleInternalServerErr(w, req)
		}

		if req.RequestLine.RequestTarget == "/" {
			return handleStatusOK(w, req)
		}

		return &server.HandlerError{
			StatusCode: response.StatusNotFound,
			Message:    "the given request target is not found",
		}
	}

	server := server.New(server.WithHandler(handler))

	defer server.Close()

	server.Listen("42069")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handleBadRequest(w *response.Writer, req *request.Request) *server.HandlerError {
	body := `
	<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
	</html>
`

	if err := w.WriteStatusLine(response.StatusBadRequest); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	headers := headers.New()
	headers.Add("Content-Length", fmt.Sprintf("%d", len(body)))
	headers.Add("Content-Type", "text/html")

	if err := w.WriteHeaders(headers); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	if _, err := w.WriteBody([]byte(body)); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return nil
}

func handleInternalServerErr(w *response.Writer, req *request.Request) *server.HandlerError {
	body := `
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>
`

	if err := w.WriteStatusLine(response.StatusInternalServerError); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	headers := headers.New()
	headers.Add("Content-Length", fmt.Sprintf("%d", len(body)))
	headers.Add("Content-Type", "text/html")

	if err := w.WriteHeaders(headers); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	if _, err := w.WriteBody([]byte(body)); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return nil
}

func handleStatusOK(w *response.Writer, req *request.Request) *server.HandlerError {
	body := `
	<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
	</html>
`

	if err := w.WriteStatusLine(response.StatusOK); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	headers := headers.New()
	headers.Add("Content-Length", fmt.Sprintf("%d", len(body)))
	headers.Add("Content-Type", "text/html")

	if err := w.WriteHeaders(headers); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	if _, err := w.WriteBody([]byte(body)); err != nil {
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return nil
}
