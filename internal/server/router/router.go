package router

import (
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/handlers"
	"github.com/Maxim-Ba/metriccollector/internal/server/handlers/middleware"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func New() *chi.Mux {
	logger.LogInfo("InitHandlers")
	r := chi.NewRouter()
	r.Get("/", middlewares(handlers.GetAllHandler))

	r.Route("/value", func(r chi.Router) {
		r.Post("/", middlewares(handlers.GetOneHandler))
		r.Get("/{metricType}/{metricName}", middlewares(handlers.GetOneHandlerByParams))
	})
	r.Route("/update", func(r chi.Router) {
		r.Post("/", middlewares(handlers.UpdateHandler))
		r.Post("/{metricType}/{metricName}/{value}", middlewares(handlers.UpdateHandlerByURLParams))
		r.Get("/{metricType}/{metricName}/{value}", middlewares(handlers.UpdateHandlerByURLParams))
	})
	return r
}

func middlewares(next http.HandlerFunc) http.HandlerFunc {
	mids := []Middleware{
		middleware.GzipHandle,
		middleware.WithLogging,
		storage.WithSyncLocalStorage,
	}
	for _, mid := range mids {
		next = mid(next)
	}
	return next
}
