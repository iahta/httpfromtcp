package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	h := make(map[string]string)
	return h
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	headerLineText := string(data[:idx])
	n = idx + 2

	colonIndex := strings.IndexByte(headerLineText, ':')
	if colonIndex == -1 {
		return 0, false, fmt.Errorf("invalid header: no colon found")
	}

	if colonIndex > 0 && headerLineText[colonIndex-1] == ' ' {
		return 0, false, fmt.Errorf("invalid header: space before colon")
	}

	key := strings.TrimSpace(headerLineText[:colonIndex])
	value := strings.TrimSpace(headerLineText[colonIndex+1:])

	if key == "" {
		return 0, false, fmt.Errorf("invalid header: empty key")
	}

	for _, c := range key {
		switch {
		case 'A' <= c && c <= 'Z':
		case 'a' <= c && c <= 'z':
		case '0' <= c && c <= '9':
		case strings.ContainsRune("!#$%&'*+-.^_|~", c):
		default:
			return 0, false, fmt.Errorf("invalid character: %s", key)
		}
	}

	lowerKey := strings.ToLower(key)

	h[lowerKey] = value

	return n, false, nil
}
