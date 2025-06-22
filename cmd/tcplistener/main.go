package main

import (
	"fmt"
	"log"
	"net"

	"github.com/iahta/httpfromtcp/internal/request"
)

func main() {
	tcpListen, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not make connection: %s\n", err)
	}
	defer func() {

		tcpListen.Close()
	}()

	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			log.Fatalf("could not accept connection: %s\n", err)
		}
		fmt.Printf("Connection has been accepted\n")

		channelLines, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("unable to parse request")
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", channelLines.RequestLine.Method, channelLines.RequestLine.RequestTarget, channelLines.RequestLine.HttpVersion)
	}

}
