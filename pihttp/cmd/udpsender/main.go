package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("error to Resolve UDP Addr err: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("error Dial UDP conn err: %s", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")

		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error on reader.ReadString err: %s", err)
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Printf("error on conn.Write err: %s", err)
		}
	}

}
