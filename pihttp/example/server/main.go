package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gpbPiazza/pinet/pihttp"
)

func main() {
	s := pihttp.NewServer()

	s.HandleFunc(http.MethodGet, "/text-time", func(req pihttp.Request, resp *pihttp.Response) error {
		resp.Body = []byte("Get Many times sending just text for get time")
		resp.Header["Content-Type"] = []string{"text"}
		return nil
	})

	s.HandleFunc(http.MethodPost, "/text-time", func(req pihttp.Request, resp *pihttp.Response) error {
		resp.Body = []byte("POST text-time sending just text for get time")
		resp.Header["Content-Type"] = []string{"text"}
		return nil
	})

	s.HandleFunc(http.MethodDelete, "/text-time", func(req pihttp.Request, resp *pihttp.Response) error {
		resp.Body = []byte("DELETE text-time sending just text for get time")
		resp.Header["Content-Type"] = []string{"text"}
		return nil
	})

	s.HandleFunc(http.MethodPost, "/time", func(req pihttp.Request, resp *pihttp.Response) error {
		// `{"nome": "vamo gremio", "eh_os_guri": true, "torcedor_maluco_ids": [1,2,3,4,5]}`
		type TextTimeBody struct {
			Nome                string `json:"Nome"`
			Eh_os_guri          bool   `json:"eh_os_guri"`
			Torcedor_maluco_ids []int  `json:"torcedor_maluco_ids"`
		}

		body := new(TextTimeBody)
		if err := json.Unmarshal(req.EntityBody, body); err != nil {
			log.Printf("error or parse body err : %s", err)
			return err
		}

		log.Printf("received body -> %t - %v - %s", body.Eh_os_guri, body.Torcedor_maluco_ids, body.Nome)

		respBody, err := json.Marshal(*body)
		if err != nil {
			log.Printf("error or parse body response err : %s", err)
			return err
		}

		resp.Body = respBody
		resp.Header["Content-Type"] = []string{"application/json"}
		resp.StatusCode = http.StatusOK
		return nil
	})

	s.Start()
}
