package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"goOrderAPI/controllers"
	"time"
)

func PopulateRoutes() *chi.Mux {

	//Middleware Stack
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	requestLogger := logrus.New()
	requestLogger.Formatter = &logrus.TextFormatter{}
	requestLogger.Level = logrus.InfoLevel
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger: requestLogger,
		},
	)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.NotFound(controllers.HandleNotFound)
	router.MethodNotAllowed(controllers.HandleMethodNotAllowed)

	router.Route("/orders", func(r chi.Router) {
		r.Post("/", controllers.CreateOrder)
		r.Patch("/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}}", controllers.TakeOrder)
		r.Get("/", controllers.ListOrders)
	})

	return router
}
