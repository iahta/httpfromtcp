package response

import (
	"fmt"
	"io"

	"github.com/iahta/httpfromtcp/internal/headers"
)

type Writer struct {
	writer      io.Writer
	writerState WriterState
}

type WriterState int

const (
	WriterStatusCode WriterState = iota
	WriterHeaders
	WriterBody
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: WriterStatusCode,
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != WriterStatusCode {
		return fmt.Errorf("cannont write status line in state %d", w.writerState)
	}
	defer func() { w.writerState = WriterHeaders }()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != WriterHeaders {
		return fmt.Errorf("cannont write status line in state %d", w.writerState)
	}
	defer func() { w.writerState = WriterBody }()
	for key, value := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != WriterBody {
		return 0, fmt.Errorf("cannont write status line in state %d", w.writerState)
	}
	return w.writer.Write(p)
}
