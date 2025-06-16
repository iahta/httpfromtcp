package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	currentLine := ""

	for {
		b := make([]byte, 8)
		n, err := file.Read(b)
		if err != nil {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
			}
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		str := string(b[:n])
		parts := strings.Split(str, "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLine, parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}

}
