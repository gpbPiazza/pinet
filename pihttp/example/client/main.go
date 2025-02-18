package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("you must provid a arg! http or tcp arg")
	}

	switch args[1] {
	case "http":
		doHTTPRequest()
	case "tcp":
		doTCPRequest()
	default:
		log.Fatal("not reconized arg passed, provided 'http' or 'tcp' args")
	}
}

func doHTTPRequest() {
	res, err := http.Get("http://localhost:8080")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Client got content length %d", res.ContentLength)
	log.Printf("Client got statusCode %d", res.StatusCode)

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
