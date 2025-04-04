package main

import (
	"github.com/gpbPiazza/httpfromtcp/internal/server"
)

func main() {
	server := server.New()

	defer server.Close()

	server.Listen("42069")

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// <-sigChan
	// log.Println("Server gracefully stopped")
}
