package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gpbPiazza/httpfromtcp/internal/request"
	"github.com/gpbPiazza/httpfromtcp/internal/response"
	"github.com/gpbPiazza/httpfromtcp/internal/server"
)

func main() {
	handler := func(w *response.Writer, req *request.Request) {
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
			handlerProxyStream(w, req)
			return
		}

		if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
			handleVideo(w, req)
			return
		}

		if req.RequestLine.RequestTarget == "/yourproblem" {
			handler400(w, req)
			return
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			handler500(w, req)
			return
		}

		if req.RequestLine.RequestTarget == "/" {
			handler200(w, req)
			return
		}

		handler404(w, req)
		return
	}

	server := server.New(server.WithHandler(handler))

	defer server.Close()

	server.Listen("42069")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handleVideo(w *response.Writer, req *request.Request) {
	file, err := os.Open("./assets/vim.mp4")
	if err != nil {
		log.Printf("err to open file err: %s", err)
		handler500(w, req)
		return
	}

	body, err := io.ReadAll(file)
	if err != nil {
		log.Printf("err to readAll file err: %s", err)
		handler500(w, req)
		return
	}

	w.WriteStatusLine(response.StatusOK)
	h := response.DefaultHeaders(len(body))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handlerProxyStream(w *response.Writer, req *request.Request) {
	requestTarget, _ := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin")
	resp, err := http.Get(fmt.Sprintf("%s/%s", "https://httpbin.org", requestTarget))
	if err != nil {
		log.Printf("error proxing request err: %s", err)
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := response.DefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Add("Trailer", "X-Content-SHA256")
	h.Add("Trailer", "X-Content-Length")
	h.Delete("Content-Length")
	w.WriteHeaders(h)

	oneKbChunk := 1024
	buf := make([]byte, oneKbChunk)
	var rawBody []byte
	for {
		numBytesRead, err := resp.Body.Read(buf)
		log.Printf("number of bytes readed from body - %d", numBytesRead)
		if errors.Is(err, io.EOF) {
			log.Print("end of file error - breaking while loop")
			break
		}
		if err != nil {
			log.Printf("err reading response body err: %s", err)
			break
		}

		if numBytesRead > 0 {
			rawBody = append(rawBody, buf[:numBytesRead]...)
			if _, err = w.WriteChunkedBody(buf[:numBytesRead]); err != nil {
				log.Printf("err WriteChunkedBody err: %s", err)
				break
			}
		}
	}

	if _, err = w.WriteChunkedBodyDone(); err != nil {
		log.Printf("err WriteChunkedBodyDone err: %s", err)
	}

	hashRawBody := sha256.Sum256(rawBody)
	h.Add("X-Content-SHA256", fmt.Sprintf("%x", hashRawBody))
	h.Add("X-Content-Length", fmt.Sprintf("%d", len(rawBody)))
	if err = w.WriteTrailers(h); err != nil {
		log.Printf("err WriteTrailers err: %s", err)
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.DefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.DefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.DefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler404(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusNotFound)
	body := []byte(`<html>
<head>
<title>404 Not found</title>
</head>
<body>
<h1>Not found!</h1>
<p>Your request suck! I dont know this targer bro!</p>
</body>
</html>
`)
	h := response.DefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
