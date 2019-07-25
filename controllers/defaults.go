package controllers

import (
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"goOrderAPI/logger"
	"net/http"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status,omitempty"` // user-level status message
	ErrorText  string `json:"error,omitempty"`  // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      err.Error(),
	}
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	err := render.Render(w, r, ErrNotFound)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Cannot Render")
	}
}

func HandleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	err := render.Render(w, r, &ErrResponse{HTTPStatusCode: http.StatusMethodNotAllowed, ErrorText: "Method Not Allowed"})
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"err": err}).Error("Cannot render")
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, ErrorText: "Resource not found."}

type SuccessResponse struct {
	HTTPStatusCode int    `json:"-"`                // http response status code
	StatusText     string `json:"status,omitempty"` // user-level status message
}

func (s *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HTTPStatusCode)
	return nil
}

var RequestSuccessfull = &SuccessResponse{HTTPStatusCode: http.StatusOK, StatusText: "SUCCESS"}
