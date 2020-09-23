package routes

import (
	"github.com/d-kolpakov/logger"
	"github.com/go-chi/chi"
	"github.com/heptiolabs/healthcheck"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"rebrain-location/internal/handlers"
	"rebrain-location/pkg/helpers/probs"
	"rebrain-location/pkg/middleware"
)

func Start(serviceName string, appVersion string, logger *logger.Logger, db *pgxpool.Pool) {
	serviceUtlPrefix := "/" + serviceName

	r := chi.NewRouter()
	m := middleware.New(logger)
	h := &handlers.Handler{Db: db, L: logger, ServiceName: serviceName, AppVersion: appVersion}

	// Health check
	hc := healthcheck.NewHandler()
	hc.AddReadinessCheck("readiness-check", func() error {
		return probs.GetReadinessErr()
	})
	hc.AddLivenessCheck("liveness-check", func() error {
		return probs.GetLivenessErr()
	})

	r.Use(m.RecoverMiddleware)
	r.Route(serviceUtlPrefix, func(r chi.Router) {
		//System endpoints
		r.Group(func(r chi.Router) {
			r.Get("/", h.HomeRouteHandler)
		})

		// Application endpoints
		r.Group(func(r chi.Router) {
			r.Use(m.ContextRequestMiddleware, m.LogRequests)
			r.Get("/endpoint/", h.InternalEndpoint)
			//Public endpoints
			r.Route("/public",func(r chi.Router) {
				r.Get("/info/", h.PublicEndpoint)
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
