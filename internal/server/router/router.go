package router

import (
	"fmt"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	fmt.Print("InitHandlers")
	r := chi.NewRouter()
	r.Get("/", handlers.GetAllHandler)

	r.Route("/value", func(r chi.Router) {
		r.Post("/", logger.WithLogging (handlers.GetOneHandler))
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/", logger.WithLogging (handlers.UpdateHandler))
		r.Post("/{metricType}/{metricName}/{value}", logger.WithLogging (handlers.UpdateHandlerByURLParams))

	})

	return r
}
