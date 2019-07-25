package main

import (
	"net/http"
	"time"
)

func NewHTTPServer() error {

	handlers := PopulateRoutes()
	server := http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handlers,
	}

	return server.ListenAndServe()
}
