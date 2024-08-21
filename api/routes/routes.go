package routes

import (
	"github.com/eniac-x-labs/manta-indexer/api/service"
	"github.com/go-chi/chi/v5"
)

type Routes struct {
	router *chi.Mux
	svc    service.Service
}

// NewRoutes ... Construct a new route handler instance
func NewRoutes(r *chi.Mux, svc service.Service) Routes {
	return Routes{
		router: r,
		svc:    svc,
	}
}
