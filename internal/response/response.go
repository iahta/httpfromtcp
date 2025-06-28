package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/iahta/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK StatusCode = iota + 1
	StatusBadRequest
	StatusInternalServerError
)

func WriteStatusLine(w io.Writer, stausCode StatusCode) error {
	switch stausCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return fmt.Errorf("error writing OK status: %v", err)
		}
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return fmt.Errorf("error writing Bad Request status: %v", err)
		}
	case StatusInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return fmt.Errorf("error writing Server Error status: %v", err)
		}
	default:
		return fmt.Errorf("unknown status code")
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

const crlf = "\r\n"

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		header := fmt.Sprintf("%s: %s%s", key, value, crlf)
		_, err := w.Write([]byte(header))
		if err != nil {
			return fmt.Errorf("error writing field-line: %v", err)
		}
	}
	_, err := w.Write([]byte(crlf))
	if err != nil {
		return fmt.Errorf("error writing CRLF: %v", err)
	}
	return nil
}
