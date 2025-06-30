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
		return 0, fmt.Errorf("cannont write body in state: %d", w.writerState)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != WriterBody {
		return 0, fmt.Errorf("cannont write body in state: %d", w.writerState)
	}
	chunkLen := []byte(fmt.Sprintf("%x\r\n", len(p)))
	_, err := w.writer.Write(chunkLen)
	if err != nil {
		return 0, fmt.Errorf("cannont write body length: %v", err)
	}
	chunkWrit, err := w.writer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("cannont write body: %v", err)
	}
	_, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, fmt.Errorf("cannont write end body: %v", err)
	}
	return chunkWrit, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != WriterBody {
		return 0, fmt.Errorf("cannont write body done in state: %d", w.writerState)
	}
	done := []byte(fmt.Sprintf("0\r\n"))
	doneLen, err := w.writer.Write(done)
	if err != nil {
		return 0, fmt.Errorf("cannont write length of end: %v", err)
	}
	_, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, fmt.Errorf("cannont write end: %v", err)
	}
	return doneLen, nil
}
