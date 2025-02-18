package main

import "github.com/gpbPiazza/pinet/pihttp"

func main() {
	s := pihttp.NewServer()

	s.Start()
}
