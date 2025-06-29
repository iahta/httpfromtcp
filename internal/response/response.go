package response

import (
	"fmt"
	"strconv"

	"github.com/iahta/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	switch statusCode {
	case StatusOK:
		return []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		return []byte(fmt.Sprintf("HTTP/1.1 400 Bad Request\r\n"))
	case StatusInternalServerError:
		return []byte(fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\n"))
	}
	return nil

}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
