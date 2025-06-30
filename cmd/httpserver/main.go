package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/iahta/httpfromtcp/internal/headers"
	"github.com/iahta/httpfromtcp/internal/request"
	"github.com/iahta/httpfromtcp/internal/response"
	"github.com/iahta/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		videoHandler(w, req)
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
		route := fmt.Sprintf("https://httpbin.org/%s", target)
		req.RequestLine.RequestTarget = route
		httpHandler(w, req)
		return

	}

	handler200(w, req)
	return
}

func videoHandler(w *response.Writer, req *request.Request) {
	data, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		handler500(w, req)
		return
	}
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(len(data))
	h.Override("Content-Type", "video/mp4")
	err = w.WriteHeaders(h)
	if err != nil {
		fmt.Printf("error writing headers: %v", err)
		handler500(w, req)
		return
	}
	_, err = w.WriteBody(data)
	if err != nil {
		fmt.Printf("error writing body: %v", err)
		handler500(w, req)
		return
	}
}

func httpHandler(w *response.Writer, req *request.Request) {
	buf := make([]byte, 32, 32)
	resp, err := http.Get(req.RequestLine.RequestTarget)
	if err != nil {
		handler400(w, req)
		return
	}
	defer resp.Body.Close()
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.OverrideContentLength()
	h.Set("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")
	w.WriteHeaders(h)

	body := []byte{}
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
			w.WriteChunkedBody(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			handler500(w, req)
			return
		}
		fmt.Printf("%v\n", n)
	}
	w.WriteChunkedBodyDone()
	sum := sha256.Sum256(body)
	t := headers.NewHeaders()
	bodyLen := strconv.Itoa(len(body))
	t.Set("X-Content-SHA256", fmt.Sprintf("%x", sum))
	t.Set("X-Content-Length", bodyLen)
	err = w.WriteTrailers(t)
	if err != nil {
		handler500(w, req)
		return
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
	</body>
	</html>
	`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := []byte(`<html>
	<head>
	<title>500 Internal Server Error</title>
	</head>
	<body>
	<h1>500 Internal Server Error</h1>
	<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>
	`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
	<head>
	<title>200 OK</title>
	</head>
	<body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
	</body>
	</html>
	`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
