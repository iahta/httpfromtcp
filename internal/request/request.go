package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/iahta/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine    RequestLine
	parserState    requestState
	Headers        headers.Headers
	Body           []byte
	bodyLengthRead int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	b := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	r := &Request{
		parserState: requestStateInitialized,
		Headers:     headers.NewHeaders(),
		Body:        make([]byte, 0),
	}
	for r.parserState != requestStateDone {
		if readToIndex >= len(b) {
			buf := make([]byte, cap(b)*2, cap(b)*2)
			copy(buf, b)
			b = buf
		}
		n, err := reader.Read(b[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.parserState != requestStateDone {
					return nil, fmt.Errorf("incomplete request: %s", err)
				}
				break
			}
			return nil, err
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
	totalBytesParsed := 0
	for r.parserState != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
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
		r.parserState = requestStateParsingHeaders
		return consumed, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("error failed to parse header: %v", err)
		}
		if done {
			r.parserState = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		contentLength, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.parserState = requestStateDone
			return len(data), nil
		}
		contentLengthNum, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("malformed Content-Length: %v", err)
		}
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)
		if r.bodyLengthRead > contentLengthNum {
			return 0, fmt.Errorf("error: body greater than content-length: body: %v content-length: %v", len(r.Body), contentLengthNum)
		}
		if r.bodyLengthRead == contentLengthNum {
			r.parserState = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
