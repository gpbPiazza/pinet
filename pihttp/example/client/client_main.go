package main

import (
	"log"
	"net"
)

func main() {
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
