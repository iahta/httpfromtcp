package server

import (
	"github.com/iahta/httpfromtcp/internal/request"
	"github.com/iahta/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)
