package api

import (
	db "server/database"
	redis "server/redis"
	ws "server/ws"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	DB    *db.Queries
	Redis *redis.Client
	Wsm   *ws.WebSocketManager
}

func (s *Server) RegisterRoute(r chi.Router) {
	// creaste users
	r.Post("/riders/register", s.createRider)
	r.Post("/drivers/register", s.createDriver)

	// login users
	r.Post("/riders/loginRider", s.LoginRider)
	r.Post("/drivers/loginDriver", s.LoginDriver)

	// rides handling
	r.Post("/rides/request", s.RequestRideHandler)
	r.Post("/drivers/{id}/location", s.updateDriverLocation)
	r.Get("/rides/{id}/status", s.GetRideStatusHandler)
	r.Get("/analytics", s.AnalyticsHandler)

}
