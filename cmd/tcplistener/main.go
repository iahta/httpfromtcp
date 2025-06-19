package main

import (
	"fmt"
	"log"
	"net"
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

		channelLines := getLinesChannel(conn)
		for channel := range channelLines {
			fmt.Printf("%s\n", channel)
		}
		fmt.Printf("Connection has been closed\n")
	}

}
