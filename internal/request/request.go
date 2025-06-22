package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	parserState requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	b := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	r := &Request{
		parserState: requestStateInitialized,
	}
	for r.parserState != requestStateDone {
		if cap(b) == readToIndex {
			buf := make([]byte, cap(b)*2, cap(b)*2)
			copy(buf, b)
			b = buf
		}
		n, err := reader.Read(b[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.parserState = requestStateDone
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		readToIndex += n

		parsed, err := r.parse(b[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(b, b[parsed:])
		readToIndex -= parsed
	}
	return r, nil
}

func parseRequestLine(request []byte) (*RequestLine, int, error) {
	idx := bytes.Index(request, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(request[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formattred request-line: %s", parts)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.parserState {
	case requestStateInitialized:
		reqLine, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("error failed to parse request: %v", err)
		}
		if consumed == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.parserState = requestStateDone
		return consumed, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
