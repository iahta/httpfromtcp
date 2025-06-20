package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading request: %v", err)
	}
	reqLine, err := parseRequestLine(string(req))
	if err != nil {
		return nil, fmt.Errorf("error failed to parse request: %v", err)
	}
	return &Request{
		RequestLine: reqLine,
	}, nil
}

func parseRequestLine(request string) (RequestLine, error) {
	parts := strings.Split(request, "\r\n")
	reqSplit := strings.Split(parts[0], " ")
	if len(reqSplit) < 3 || len(reqSplit) > 3 {
		return RequestLine{}, fmt.Errorf("request line must be method, target, version")
	}
	method := reqSplit[0]
	requestTarget := reqSplit[1]
	httpVersion := strings.Split(reqSplit[2], "/")
	if !IsAllUpper(method) {
		return RequestLine{}, fmt.Errorf("method must be all capital")
	}
	if httpVersion[1] != "1.1" {
		return RequestLine{}, fmt.Errorf("http version must be 1.1")
	}
	if !strings.HasPrefix(requestTarget, "/") {
		return RequestLine{}, fmt.Errorf("request target missing '/'")
	}
	requestLine := RequestLine{
		HttpVersion:   httpVersion[1],
		RequestTarget: requestTarget,
		Method:        method,
	}
	return requestLine, nil
}

//request-line  = method SP request-target SP HTTP-version
//				GET       /coffee           HTTP/1.1

func IsAllUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false // Found a non-uppercase character
		}
	}
	return true // All characters were uppercase
}
