package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChannel := make(chan string)
	go readLines(f, lineChannel)
	return lineChannel
}

func readLines(f io.ReadCloser, lineChannel chan string) {
	currentLine := ""
	for {
		b := make([]byte, 8)
		n, err := f.Read(b)
		if err != nil {
			if currentLine != "" {
				lineChannel <- currentLine
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
			currentLine += parts[i]
			lineChannel <- currentLine
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
	close(lineChannel)
	f.Close()
}
