package main

import (
	"net/http"

	"github.com/gpbPiazza/pinet/pihttp"
)

func main() {
	s := pihttp.NewServer()

	s.HandleFunc(http.MethodGet, "/time", func(req pihttp.Request, resp *pihttp.Response) error {

		return nil
	})

	s.Start()
}
