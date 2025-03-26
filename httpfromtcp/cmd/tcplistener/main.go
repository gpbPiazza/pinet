package main

import "github.com/gpbPiazza/httpfromtcp"

func main() {
	s := httpfromtcp.NewServer()
	s.Listen("42069")
}
