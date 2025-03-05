package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("you must provid a arg! http or tcp arg")
	}

	switch args[1] {
	case "httpGet":
		doGetHTTPRequest()
	case "httpPost":
		doPostHTTPRequest()
	case "tcp":
		doTCPRequest()
	default:
		log.Fatal("not reconized arg passed, provided 'http' or 'tcp' args")
	}
}

func doGetHTTPRequest() {
	urlVal := make(url.Values)

	urlVal.Add("tem_jogo_hoje", "true")
	urlVal.Add("estadios_ids", "1")
	urlVal.Add("estadios_ids", "2")
	urlVal.Add("estadios_ids", "3")
	urlVal.Add("data_jogo", "02/05/2025")

	endpoint := "http://localhost:8080" + "?" + urlVal.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Fatalf("error on new request ctor err: %s", err)
	}

	req.Header.Add("Times-Do-RJ", "flu")
	req.Header.Add("Times-Do-RJ", "fla")
	req.Header.Add("Times-Do-RJ", "fogão")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error on sending http request err: %s", err)
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Client error on res.Body.Read err: %s", err)
	}

	log.Printf("Client got content length %d", res.ContentLength)
	log.Printf("Client got statusCode %d", res.StatusCode)

	for headerKey, headerVal := range res.Header {
		for _, val := range headerVal {
			log.Printf("Client got Header %s: %s", headerKey, val)
		}
	}

	log.Printf("Client got body %s", string(resBodyBytes))
}

func doPostHTTPRequest() {
	urlVal := make(url.Values)

	urlVal.Add("tem_jogo_hoje", "true")
	urlVal.Add("estadios_ids", "1")
	urlVal.Add("estadios_ids", "2")
	urlVal.Add("estadios_ids", "3")
	urlVal.Add("data_jogo", "02/05/2025")

	resource := "time"
	host := "http://localhost:8080"
	url := host + "/" + resource + "?" + urlVal.Encode()

	body := `{"nome": "vamo gremio", "eh_os_guri": true, "torcedor_maluco_ids": [1,2,3,4,5]}`

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("error jsonMarshal err: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Fatalf("error on new request ctor err: %s", err)
	}

	req.Header.Add("Times-Do-RJ", "flu")
	req.Header.Add("Times-Do-RJ", "fla")
	req.Header.Add("Times-Do-RJ", "fogão")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error on sending http request err: %s", err)
	}
	defer res.Body.Close()

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Client error on res.Body.Read err: %s", err)
	}

	log.Printf("Client got content length %d", res.ContentLength)
	log.Printf("Client got statusCode %d", res.StatusCode)

	for headerKey, headerVal := range res.Header {
		for _, val := range headerVal {
			log.Printf("Client got Header %s: %s", headerKey, val)
		}
	}

	log.Printf("Client got body %s", string(resBodyBytes))
}

func doTCPRequest() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error connecting: %s", err)
	}
	defer conn.Close()

	request := "E o que isso a feta o gremio??"
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("Error writing to server: %s", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalf("Error reading from server: %s", err)
		return
	}

	log.Println("Server response:", string(buffer[:n]))
}
