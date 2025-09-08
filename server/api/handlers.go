package api

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	DB    *db.Queries
	Redis *redis.Client
	WSHub *ws.Hub
}

func (s *Server) RegisterRoute(r chi.Router) {
	r.Post("/riders/register", s.createRiders)
	r.Post("/drivers/register", s.createDriver)
	r.Post("/rides/request", s.createDriver)
	r.Post("/drivers/{id}/location", updateDriverLocation)
	r.GET("/rides/{id}/status", s.getRideStatus)
}

func (s *Server) createRiders(w http.ResponseWriter, r *http.Request) {
	type req struct {
		username string `json:"username"`
	}
	var body req
	json.NewDecoder(r.Body).Decode(&body)

	w.WriteHeader(http.StatusCreated)
}
