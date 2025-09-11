package api

import (
	db "uber-like-system/server/database"
	redis "uber-like-system/server/redis"
	ws "uber-like-system/server/ws"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	DB    *db.Queries
	Redis *redis.Client
	Wsm   *ws.WebSocketManager
}

func (s *Server) RegisterRoute(r chi.Router) {
	// creaste users
	r.Post("/riders/signup", s.createRider)
	r.Post("/drivers/signup", s.createDriver)

	// login users
	r.Post("/riders/login", s.LoginRider)
	r.Post("/drivers/login", s.LoginDriver)

	// rides handling
	r.Post("/rides/request", s.RequestRideHandler)
	r.Post("/riders/{1}/accept", s.AcceptRideHandler)

	r.Post("/drivers/{id}/location", s.updateDriverLocation)
	r.Get("/rides/{id}/status", s.GetRideStatusHandler)

	// admin testing
	r.Get("/drivers/nearby", s.GetNearbyDriversHandler)
	r.Get("/analytics", s.AnalyticsHandler)

}
