package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "server/database"
	redis "server/redis"

	"github.com/go-chi/chi"
)

type Server struct {
	DB    *database.Queries
	Redis *redis.Client
	WSHub *ws.Hub
}

func (s *Server) RegisterRoute(r chi.Router) {
	r.Post("/riders/register", s.createRiders)
	r.Post("/drivers/register", s.createDriver)
	// r.Post("/rides/request", s.createDriverRequest)
	r.Post("/drivers/{id}/location", s.updateDriverLocation)
	r.Get("/rides/{id}/status", s.getRideStatus)
}

func (s *Server) createRiders(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) createDriver(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Username string `json:"username"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	_, err := s.DB.Exec("Insert into drives (name) values ()", body.Username)
	if err != nil {
		http.Error(w, "faailed to create driver", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) updateDriverLocation(w http.ResponseWriter, r *http.Request) {
	type req struct {
		DriverID  int64   `json:"driver_id"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"long"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := s.Redis.HSet(r.Context(),
		fmt.Sprintf("%d", body.DriverID),
		"lat", body.Latitude,
		"long", body.Longitude,
	).Err()
	if err != nil {
		http.Error(w, "failed to update location", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getRideStatus(w http.ResponseWriter, r *http.Request) {
	rideID := chi.URLParam(r, "id")

	var status string
	err := s.DB.QueryRow("SELECT status FROM ride_requests WHERE id=$1", rideID).Scan(&status)
	if err != nil {
		http.Error(w, "ride not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
