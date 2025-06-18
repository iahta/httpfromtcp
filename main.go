package main

import (
	"fmt"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	channelLines := getLinesChannel(file)
	for channel := range channelLines {
		fmt.Printf("read: %s\n", channel)
	}
}
