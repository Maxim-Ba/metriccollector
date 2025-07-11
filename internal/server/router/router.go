package router

import (
	"net/http"

	m "github.com/go-chi/chi/v5/middleware"

	"github.com/Maxim-Ba/metriccollector/internal/server/handlers"
	"github.com/Maxim-Ba/metriccollector/internal/server/handlers/middleware"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

// Middleware is a type alias for functions that wrap http.HandlerFunc
// to provide additional functionality like logging, compression, etc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// New creates and configures a new chi.Mux router with all application routes
// and middleware. The router includes:
// - Debug profiling endpoints under /debug
// - Metric retrieval and update endpoints
// - Database health check endpoint
// Middlewares are applied in the order: signature verification, storage sync,
// gzip compression, and request logging.
func New() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/debug", m.Profiler())

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
	r.Route("/updates", func(r chi.Router) {
		r.Post("/", middlewares(handlers.UpdatesHandler))
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", middlewares(handlers.PingDB))
	})
	return r
}

func middlewares(next http.HandlerFunc) http.HandlerFunc {
	mids := []Middleware{
		middleware.SignatureHandle,
		storage.WithSyncLocalStorage,
		middleware.GzipHandle,
		middleware.WithLogging,
	}
	for _, mid := range mids {
		next = mid(next)
	}
	return next
}
