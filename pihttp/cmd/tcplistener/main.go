package main

import "github.com/gpbPiazza/httpfromtcp/internal/server"

func main() {
	s := server.New()
	s.Listen("42069")
}
