package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	network, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("error resolving connection: %v", err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, network)
	if err != nil {
		log.Fatalf("error dialing UDP: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)

	for {

		fmt.Println(">")
		str, err := r.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
			os.Exit(1)
		}
		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatalf("Error sending message %v", err)
			os.Exit(1)
		}
		fmt.Printf("Message sent: %s", str)
	}

}
