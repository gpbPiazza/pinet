package main

import "github.com/gpbPiazza/pinet/httpfromtcp"

func main() {
	s := httpfromtcp.NewServer()
	s.Listen("42069")
}
